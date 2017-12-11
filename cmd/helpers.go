// Copyright © 2017 Cisco Systems, Inc.
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

package cmd

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/codeskyblue/go-sh"
	"github.com/smallfish/simpleyaml"
	"github.com/spf13/viper"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/engine"
	"k8s.io/helm/pkg/proto/hapi/chart"

	log "github.com/sirupsen/logrus"
)

func buildDependencies(app string, chart string) {
	// If requirements.yaml exists at chart, build dependencies for it,
	// in the chart's own directory.
	if _, err := os.Stat(chart + "/" + "requirements.yaml"); !os.IsNotExist(err) {
		log.WithFields(log.Fields{
			"app":   app,
			"chart": chart,
			"err":   err,
		}).Info("Building chart dependencies for app.")
		// We start a new shell session here so
		// that we needn't tangle with `cd`.
		session := sh.NewSession()
		session.SetDir(chart)
		out, err := session.Command(
			"helm", "dependency", "build",
		).Output()
		if err != nil {
			log.WithFields(log.Fields{
				"app":   app,
				"chart": chart,
				"err":   err,
				"out":   out,
			}).Fatal("Failed to build chart dependencies for app.")
		}
		session.ShowCMD = true
	}
}

func getCurrentContext() string {
	cmd := []interface{}{
		"config",
		"current-context",
	}

	out, err := sh.Command("kubectl", cmd...).Output()
	if err != nil {
		log.Fatal("Failed running `kubectl config current-context`.")
	}

	return strings.TrimSuffix(string(out), "\n")
}

func lateInit(cmd string) {
	targetContext := viper.GetString("targetContext")

	if targetContext != currentContext {
		log.WithFields(log.Fields{
			"tryCfgFile":     viper.ConfigFileUsed(),
			"currentContext": currentContext,
			"targetContext":  targetContext,
		}).Fatal("`kubectl config current-context` does not match config's `targetContext`.")
	}

	log.WithFields(log.Fields{
		"appsPath":       viper.GetString("appsPath"),
		"tryCfgFile":     viper.ConfigFileUsed(),
		"currentContext": currentContext,
		"targetContext":  targetContext,
		"versionNumber":  versionNumber,
	}).Info("Initializing MultiHelm.")

	if cmd != "version" {
		log.Infof("`%s` called.", cmd)
	}

	return
}

func render(app string) (string, []byte, error) {
	log.WithFields(log.Fields{
		"app": app,
	}).Info("Rendering app.")

	configFile := viper.ConfigFileUsed()
	config, err := chartutil.ReadValuesFile(configFile)
	if err != nil {
		log.WithFields(log.Fields{
			"app": app,
		}).Fatal("Failed to load values while rendering app file.")
	}

	appFile := viper.GetString("appsPath") + "/" + app + ".yaml"

	data, err := ioutil.ReadFile(appFile)
	if err != nil {
		log.WithFields(log.Fields{
			"app":     app,
			"appFile": appFile,
		}).Fatal("Failed to load app template while rendering app file.")
	}

	fakeChart := &chart.Chart{
		Metadata: &chart.Metadata{
			Name:    "fake",
			Version: "0.1.",
		},
		Templates: []*chart.Template{
			{Name: "templates/main", Data: data},
		},
	}

	out, err := engine.New().Render(fakeChart, config)
	if err != nil {
		log.WithFields(log.Fields{
			"app": app,
			"err": err,
		}).Fatal("Failed to render app file.")
	}

	overrideValues := []byte(out["fake/templates/main"])

	yml, err := simpleyaml.NewYaml(overrideValues)
	if err != nil {
		log.WithFields(log.Fields{
			"app": app,
			"err": err,
		}).Fatal("Loading of override YAML failed for app.")
	}

	chart, err := yml.Get("chart").String()
	if err != nil {
		log.WithFields(log.Fields{
			"app": app,
			"err": err,
		}).Fatal("Lookup of chart in override YAML failed for app.")
	}

	buildDependencies(app, chart)

	return chart, overrideValues, nil
}
