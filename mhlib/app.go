// Copyright © 2018 Cisco Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mhlib

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"text/template"

	"github.com/imdario/mergo"
	"github.com/sirupsen/logrus"

	"github.com/codeskyblue/go-sh"
	"github.com/ghodss/yaml"
	"github.com/hairyhenderson/gomplate"
	"github.com/hairyhenderson/gomplate/data"
	"github.com/smallfish/simpleyaml"
	"github.com/stoewer/go-strcase"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/engine"
	"k8s.io/helm/pkg/proto/hapi/chart"
	"k8s.io/helm/pkg/strvals"

	log "github.com/sirupsen/logrus"
)

// AppConfig is what can be defined in a mh configuration file and is used to
// create an App struct. It is a superset of MHConfig to enable app-specific
// configuration overrides of all mh configuration settings.
//
// Maybe: Get rid of Alias in favor of ID
type AppConfig struct {
	Alias     string   `yaml:"alias"`
	File      *AppFile `yaml:"file"`
	Key       string   `yaml:"key"`
	Name      string   `yaml:"name"`
	Namespace string   `yaml:"namespace"`
	MHConfig
}

// AppConfigs is an array of AppConfig as defined in a mh configuration file.
type AppConfigs []AppConfig

// App contains attributes defining a mh app to run and app-specific mh
// configuration overrides.
type App struct {
	AppConfig
	ID  string
	log *logrus.Entry
}

// NewApp returns an App based on a appConfig and global MHConfig defaults.
func NewApp(logger *logrus.Entry, appConfig AppConfig, mhConfig MHConfig) (*App, error) {
	// Sanitize configuration
	// Todo: sanitize more, no dashes allowed etc.
	if appConfig.Name == "" {
		return nil, fmt.Errorf("Empty name for app: %v", appConfig)
	}

	// Set configuration defaults if not overridden
	err := mergo.Merge(&appConfig.MHConfig, mhConfig)
	if err != nil {
		return nil, err
	}

	// Set App ID, prioritize Alias over Name.
	var id string
	if appConfig.Alias != "" {
		id = appConfig.Alias
	} else {
		id = appConfig.Name
	}

	// Set App Key to default of ".ID" if not defined
	if appConfig.Key == "" {
		appConfig.Key = fmt.Sprintf(".%s", strcase.LowerCamelCase(id))
	}

	return &App{
		appConfig,
		id,
		logger.WithField("app", appConfig.Name),
	}, nil
}

// Build app's chart dependencies.
//
// If requirements.yaml exists at app's chart, run `helm dependency build`
// to build dependencies at that chart's directory.
func (a *App) Build(chart string) error {
	requirementsFile := chart + "/" + "requirements.yaml"
	if _, err := os.Stat(requirementsFile); !os.IsNotExist(err) {
		a.log.WithFields(log.Fields{
			"chart":            chart,
			"requirementsFile": requirementsFile,
		}).Info("Building chart dependencies for app.")

		// Start a new shell session here to avoid running `cd`.
		session := sh.NewSession()
		session.SetDir(chart)

		// Run `helm dependency build` on the chart.
		out, err := session.Command("helm", "dependency", "update").Output()
		if err != nil {
			return fmt.Errorf("Failed to build chart dependencies for app: %v", out)
		}

		session.ShowCMD = true
	}

	return nil
}

func (a *App) Destroy(purge bool) (*[]interface{}, error) {
	a.log.Info("Destroying app")
	cmd := []interface{}{"delete", a.ID}
	if purge {
		cmd = append(cmd, "--purge")
	}
	err := sh.Command("helm", cmd...).Run()
	if err != nil {
		return &cmd, fmt.Errorf("Helm delete failed for app")
	}

	return nil, nil
}

func (a *App) Status() error {
	err := sh.Command("helm", "status", a.ID).Run()
	if err != nil {
		return fmt.Errorf("Helm status failed")
	}

	return nil
}

func (a *App) Apply(configFile string) (*[]interface{}, error) {
	a.log.Info("Applying app")
	return a.apply(configFile, false)
}

func (a *App) apply(configFile string, simulate bool) (*[]interface{}, error) {
	chart, chartVersion, overrides, err := a.render(configFile)
	if err != nil {
		return nil, err
	}

	if a.PrintRendered {
		fmt.Print(string(*overrides))
	}

	// Prepare to do `helm upgrade`
	cmd := []interface{}{"upgrade", a.ID, *chart}

	// "specify the exact chart version to install. If this is not specified, the latest version is installed"
	if chartVersion != nil {
		cmd = append(cmd, "--version", *chartVersion)
	}

	if simulate {
		// "enable verbose output"
		cmd = append(cmd, "--debug")

		// "simulate an upgrade"
		cmd = append(cmd, "--dry-run")
	}

	// "force resource update through delete/recreate if needed"
	cmd = append(cmd, "--force")

	// "if a release by this name doesn't already exist, run an install"
	cmd = append(cmd, "--install")

	if !a.NoRecreatePods {
		// "performs pods restart for the resource if applicable"
		cmd = append(cmd, "--recreate-pods")
	}

	// "namespace to install the release into. (Defaults to helm default behaviour => kubeconfig checked out ns)"
	if a.Namespace != "" {
		cmd = append(cmd, "--namespace", a.Namespace)
	}

	// Make `helm upgrade` read overrides from stdin
	cmd = append(cmd, "--values", "-")

	// Run `helm upgrade
	err = sh.Command("helm", cmd...).SetInput(string(*overrides)).Run()
	if err != nil {
		return &cmd, err
	}
	return &cmd, nil
}

func (a *App) Simulate(configFile string) (*[]interface{}, error) {
	a.log.Info("Simulating app")
	return a.apply(configFile, true)
}

func (a *App) render(configFile string) (*string, *string, *[]byte, error) {
	var chartVersion string

	// read the mh main.yaml
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Failed to read configFile %v: %v", configFile, err)
	}

	// Self-render the main.yaml with gomplate functions and datasources
	//   This does not apply to the app.yaml files.
	contents := string(data)
	renderedContents, err := selfRender(contents)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Failed to selfRender configFile %v: %v", configFile, err)
	}

	config, err := chartutil.ReadValues([]byte(renderedContents))
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Failed to load values from configFile: %v", err)
	}

	appData, err := ioutil.ReadFile(*a.File.Path) // app.yaml
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Failed to load data from appFile: %v", err)
	}

	// creating a literal
	data = []byte(
		"{{- $name := \"" + a.ID + "\" }}\n" + "{{- $app := " + a.Key + " }}\n",
	)

	// combining literal with app.yaml
	data = append(data, appData...)

	// Fakechart to send to helm rendering engine
	fakeChart := &chart.Chart{
		Metadata: &chart.Metadata{
			Name:    "fake",
			Version: "0.1.0",
		},
		Templates: []*chart.Template{
			{Name: "templates/main", Data: data},
		},
	}

	// Add config via --set command
	for _, value := range a.MHConfig.SETValues {
		err := strvals.ParseInto(value, config)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("Failed to parse values provided via --set : %v", err)
		}
	}

	out, err := engine.New().Render(fakeChart, config)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Helm rendering engine failed to render fakeChart: %v", err)
	}

	overrides := []byte(out["fake/templates/main"])

	yml, err := simpleyaml.NewYaml(overrides)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Failed to load newly rendered overrides YAML: %v", err)
	}

	chart, err := yml.Get("chart").String()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Failed to lookup chart in overrides YAML: %v", err)
	}
	chartVersion, _ = yml.Get("version").String()

	// If key-value "chart" inside app YAML is determined to be a file path,
	// build/update dependencies for it. If not a path, we needn't build
	// for it.
	re := regexp.MustCompile("(^(\\.).*)|(^/.*)")
	isPath := re.MatchString(chart)
	if isPath {
		err = a.Build(chart)
		if err != nil {
			return nil, nil, nil, err
		}
	}

	return &chart, &chartVersion, &overrides, nil
}

func selfRender(templateValuesStr string) (string, error) {
	/*
		This function will accept an input string and run it through the
		  templating engine as both the values dictionary as well as the
		  template string.  It will repeat the process until the templating
		  result stops changing.
	*/

	type gomplateConfig struct {
		Gomplate gomplate.Config `yaml:"gomplate,omitempty"`
	}
	var gomp gomplateConfig
	err := yaml.Unmarshal([]byte(templateValuesStr), &gomp)
	if err != nil {
		return "", err
	}

	funcs := template.FuncMap{
		//Add case insensitive handling as marshalled inline structures are public by default, but may be lowercase templates
		"MyEq": strings.EqualFold,
	}
	if len(gomp.Gomplate.DataSources) > 0 {
		d, err := data.NewData(gomp.Gomplate.DataSources, gomp.Gomplate.DataSourceHeaders)
		if err != nil {
			return "", err
		}
		for k, v := range gomplate.Funcs(d) {
			funcs[k] = v
		}
	}
	tmpl := template.New("SelfTemplate").
		Delims("[[", "]]").
		Option("missingkey=error").
		Funcs(funcs)

	lastRender := templateValuesStr
	for i := 0; i < 10; i++ {

		// Unmarshal the file as a values dict
		vals := map[string]interface{}{}
		err := yaml.Unmarshal([]byte(lastRender), &vals)
		if err != nil {
			return "", err
		}

		// Run the the file through the tempating engine as both values
		//   file and template file
		tmpl.Parse(string(lastRender))
		if err != nil {
			return "", err
		}
		out := new(bytes.Buffer)
		err = tmpl.Execute(out, vals)
		if err != nil {
			return "", err
		}

		newRender := out.String()
		if lastRender == newRender {
			return newRender, nil // self-templating succeeded
		} else {
			lastRender = newRender
		}
	}

	return lastRender, errors.New("Self-templating failed")
}
