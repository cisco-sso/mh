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

package cmd

import (
	"fmt"
	"strings"

	"github.com/codeskyblue/go-sh"
	"github.com/spf13/viper"

	lib "***REMOVED***/***REMOVED***/multihelm/multihelmlib"
	log "github.com/sirupsen/logrus"
)

func getApps(args []string) *lib.Apps {
	var (
		apps       []lib.App
		configApps []lib.App
		foundApp   bool
	)
	configApps = getConfigApps()
	if len(args) > 0 {
		for _, arg := range args {
			foundApp = false
			for _, configApp := range configApps {
				if arg == configApp.Alias || arg == configApp.Name {
					apps = append(apps, configApp)
					foundApp = true
					break
				}
			}
			if foundApp == false {
				log.WithFields(log.Fields{
					"apps":       apps,
					"arg":        arg,
					"configApps": configApps,
					"foundApp":   foundApp,
				}).Fatal("Command line app name/alias not found in config.")
			}
		}
	} else {
		apps = configApps
	}
	return &lib.Apps{
		Apps: apps,
	}
}

func getAppsPath() string {
	return viper.GetString("appsPath")
}

func getConfigApps() []lib.App {
	var apps []lib.App
	err := viper.UnmarshalKey("apps", &apps)
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Failed to unmarshal 'apps:' from config.")
	}
	return apps
}

func getConfigFile() string {
	return viper.ConfigFileUsed()
}

func getCurrentContext() string {
	out, err := sh.Command("kubectl", "config", "current-context").Output()
	if err != nil {
		log.Fatal("Failed running `kubectl config current-context`.")
	}
	currentContext := strings.TrimSuffix(string(out), "\n")
	return currentContext
}

func getPurge() bool {
	return viper.GetBool("purge")
}

func getPrintRendered() bool {
	return viper.GetBool("printRendered")
}

func getTargetContext() string {
	return viper.GetString("targetContext")
}

func lateInit(cmd string) {
	configFile := getConfigFile()
	currentContext := getCurrentContext()
	targetContext := getTargetContext()

	if targetContext != currentContext {
		log.WithFields(log.Fields{
			"configFile":     configFile,
			"currentContext": currentContext,
			"targetContext":  targetContext,
		}).Fatal("`kubectl config current-context` does not match config's `targetContext`.")
	}

	log.WithFields(log.Fields{
		"appsPath":       getAppsPath(),
		"cmd":            cmd,
		"configFile":     configFile,
		"currentContext": currentContext,
		"targetContext":  targetContext,
		"versionNumber":  versionNumber,
	}).Info("Initializing MultiHelm.")
}

func logVersion() {
	log.Info("MultiHelm " + versionNumber)
}

func printLicense() {
	fmt.Println(`Copyright © 2018 Cisco Systems, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.`)
}
