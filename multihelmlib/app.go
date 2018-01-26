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

package multihelmlib

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"

	"github.com/codeskyblue/go-sh"
	"github.com/smallfish/simpleyaml"
	"github.com/stoewer/go-strcase"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/engine"
	"k8s.io/helm/pkg/proto/hapi/chart"

	log "github.com/sirupsen/logrus"
)

type App struct {
	Alias string
	File  string
	Key   string
	Name  string
}

func (a *App) Apply(configFile string, appSources []AppSource,
	printRendered bool) {
	method := "apply"
	simulate := false
	cmd, err := a.apply(configFile, appSources, printRendered, simulate)
	if err != nil {
		appLog := &AppLog{
			app:           a,
			appSources:    appSources,
			cmd:           cmd,
			configFile:    configFile,
			err:           err,
			method:        method,
			printRendered: printRendered,
			simulate:      simulate,
		}
		appLog.Error()
	}
}

// Build app's chart dependencies.
//
// If requirements.yaml exists at app's chart, run `helm dependency build`
// to build dependencies at that chart's directory.
func (a *App) Build(chart string) {
	requirementsFile := chart + "/" + "requirements.yaml"
	if _, err := os.Stat(requirementsFile); !os.IsNotExist(err) {
		log.WithFields(log.Fields{
			"app":              a,
			"chart":            chart,
			"requirementsFile": requirementsFile,
		}).Info("Building chart dependencies for app.")

		// Start a new shell session here to avoid running `cd`.
		session := sh.NewSession()
		session.SetDir(chart)

		// Run `helm dependency build` on the chart.
		out, err := session.Command("helm", "dependency", "update").Output()
		if err != nil {
			log.WithFields(log.Fields{
				"app":              a,
				"chart":            chart,
				"err":              err,
				"out":              out,
				"requirementsFile": requirementsFile,
			}).Fatal("Failed to build chart dependencies for app.")
		}
		session.ShowCMD = true
	}
}

func (a *App) Destroy(purge bool) {
	method := "destroy"
	cmd := []interface{}{"delete", a.Id()}
	if purge {
		cmd = append(cmd, []interface{}{"--purge"}...)
	}
	err := sh.Command("helm", cmd...).Run()
	if err != nil {
		appLog := &AppLog{
			app:    a,
			cmd:    cmd,
			err:    err,
			method: method,
			purge:  purge,
		}
		appLog.Info("Helm delete failed for app. Continuing anyway.")
	}
}

func (a *App) GetFile(appSources []AppSource) string {
	var appFile, possibleFile string
	method := "getAppFile"
	if a.File != "" {
		if _, err := os.Stat(a.File); os.IsNotExist(err) {
			appLog := &AppLog{
				app:    a,
				method: method,
				reason: "App.File not found. App.File overrides []AppSource lookup.",
			}
			appLog.Error()
			return ""

		}
		return a.File
	}
	for _, appSource := range appSources {
		if appSource.Kind == "path" {
			possibleFile = appSource.Source + "/" + a.Name + ".yaml"
		} else {
			appLog := &AppLog{
				app:        a,
				appSources: appSources,
				method:     method,
				reason:     "One or more []AppSource has an unsupported Kind. Supported Kind values are: path",
			}
			appLog.Error()
			return ""
		}
		if _, err := os.Stat(possibleFile); os.IsNotExist(err) {
			continue
		}
		appFile = possibleFile
		break
	}
	if appFile == "" {
		appLog := &AppLog{
			app:        a,
			appSources: appSources,
			method:     method,
			reason:     "Failed to locate app file.",
		}
		appLog.Error()
		return ""
	}
	return appFile
}

// Return app.Alias if one is set.
// If a.Alias is not set, return a.Name.
func (a *App) Id() string {
	var id string
	method := "id"
	if a.Alias != "" {
		id = a.Alias
	} else if a.Name != "" {
		id = a.Name
	} else {
		appLog := &AppLog{
			app:    a,
			id:     id,
			method: method,
			reason: "Neither 'app.Alias' nor 'app.Name' were found.",
		}
		appLog.Error()
	}
	return id
}

func (a *App) GetKey() string {
	method := "getkey"
	if a.Key != "" {
		return a.Key
	}
	id := a.Id()
	if id != "" {
		return "." + strcase.LowerCamelCase(id)
	}
	appLog := &AppLog{
		app:    a,
		id:     id,
		method: method,
		reason: "Failed to determine app key. Please consider defining 'key:' on your app in your MultiHelm config.",
	}
	appLog.Error()
	return ""
}

func (a *App) Simulate(configFile string, appSources []AppSource, printRendered bool) {
	method := "simulate"
	simulate := true
	cmd, err := a.apply(configFile, appSources, printRendered, simulate)
	if err != nil {
		appLog := &AppLog{
			app:           a,
			appSources:    appSources,
			cmd:           cmd,
			configFile:    configFile,
			err:           err,
			method:        method,
			printRendered: printRendered,
			simulate:      simulate,
		}
		appLog.Error()
	}
}

func (a *App) Status() {
	method := "status"
	cmd := []interface{}{"status", a.Id()}
	err := sh.Command("helm", cmd...).Run()
	if err != nil {
		appLog := &AppLog{
			app:    a,
			cmd:    cmd,
			err:    err,
			method: method,
		}
		appLog.Info("Helm status failed for app. Continuing anyway.")
	}
}

func (app *App) apply(configFile string, appSources []AppSource,
	printRendered bool, simulate bool) ([]interface{}, error) {

	chart, chartVersion, overrides, err := app.render(configFile, appSources)
	if err != nil {
		return nil, err
	}

	if printRendered {
		fmt.Print(string(overrides))
	}

	// Prepare to do `helm upgrade`
	cmd := []interface{}{"upgrade", app.Id(), chart}

	// "specify the exact chart version to install. If this is not specified, the latest version is installed"
	if chartVersion != "" {
		cmd = append(cmd, []interface{}{"--version", chartVersion}...)
	}

	if simulate {
		// "enable verbose output"
		cmd = append(cmd, []interface{}{"--debug"}...)

		// "simulate an upgrade"
		cmd = append(cmd, []interface{}{"--dry-run"}...)
	}

	// "force resource update through delete/recreate if needed"
	cmd = append(cmd, []interface{}{"--force"}...)

	// "if a release by this name doesn't already exist, run an install"
	cmd = append(cmd, []interface{}{"--install"}...)

	// "performs pods restart for the resource if applicable"
	cmd = append(cmd, []interface{}{"--recreate-pods"}...)

	// Make `helm upgrade` read overrides from stdin
	cmd = append(cmd, []interface{}{"--values", "-"}...)

	// Run `helm upgrade
	err = sh.Command("helm", cmd...).SetInput(string(overrides)).Run()
	if err != nil {
		return cmd, err
	}
	return cmd, nil
}

func (a *App) render(configFile string, appSources []AppSource) (string, string, []byte, error) {
	var chartVersion string
	method := "render"
	appLog := &AppLog{
		app:        a,
		appSources: appSources,
		configFile: configFile,
		method:     method,
	}
	appLog.Info("Running '" + method + "' for app '" + a.Id() + "'")

	config, err := chartutil.ReadValuesFile(configFile)
	if err != nil {
		appLog := &AppLog{
			app:        a,
			appSources: appSources,
			configFile: configFile,
			err:        err,
			method:     method,
			reason:     "Failed to load values from configFile.",
		}
		appLog.Error()
	}

	appFile := a.GetFile(appSources)
	appData, err := ioutil.ReadFile(appFile)
	if err != nil {
		appLog := &AppLog{
			app:        a,
			appFile:    appFile,
			appSources: appSources,
			configFile: configFile,
			err:        err,
			method:     method,
			reason:     "Failed to load data from appFile.",
		}
		appLog.Error()
	}

	data := []byte(
		"{{- $name := \"" + a.Id() + "\" }}\n" + "{{- $app := " + a.GetKey() + " }}\n",
	)

	data = append(data, appData...)

	fakeChart := &chart.Chart{
		Metadata: &chart.Metadata{
			Name:    "fake",
			Version: "0.1.0",
		},
		Templates: []*chart.Template{
			{Name: "templates/main", Data: data},
		},
	}

	out, err := engine.New().Render(fakeChart, config)
	if err != nil {
		appLog := &AppLog{
			app:        a,
			appFile:    appFile,
			appSources: appSources,
			configFile: configFile,
			data:       data,
			err:        err,
			method:     method,
			reason:     "Helm rendering engine failed to render fakeChart.",
		}
		appLog.Error()
	}

	overrides := []byte(out["fake/templates/main"])

	yml, err := simpleyaml.NewYaml(overrides)
	if err != nil {
		appLog := &AppLog{
			app:        a,
			appFile:    appFile,
			appSources: appSources,
			configFile: configFile,
			err:        err,
			method:     method,
			reason:     "Failed to load newly rendered overrides YAML.",
		}
		appLog.Error()
	}

	chart, err := yml.Get("chart").String()
	if err != nil {
		appLog := &AppLog{
			app:        a,
			appFile:    appFile,
			appSources: appSources,
			configFile: configFile,
			err:        err,
			method:     method,
			reason:     "Failed to lookup chart in overrides YAML.",
		}
		appLog.Error()
	}
	chartVersion, _ = yml.Get("version").String()

	// If key-value "chart" inside app YAML is determined to be a file path,
	// build/update dependencies for it. If not a path, we needn't build
	// for it.
	re := regexp.MustCompile("(^(\\.).*)|(^/.*)")
	isPath := re.MatchString(chart)
	if isPath {
		a.Build(chart)
	}

	return chart, chartVersion, overrides, nil
}
