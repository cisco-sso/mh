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
	"github.com/codeskyblue/go-sh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

// simulateCmd represents the simulate command
var simulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Simulate apps",
	Long: `Simulate the apply of one or more MultiHelm apps. If you do not specify one or more
apps, MultiHelm acts on all apps in your MultiHelm config.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("Simulate called.")
		if len(args) > 0 {
			for _, arg := range args {
				simulate(arg)
			}
		} else {
			for _, arg := range viper.GetStringSlice("apps") {
				simulate(arg)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(simulateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// simulateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// simulateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func simulate(app string) {

	chart, overrideValues, err := render(app)
	if err != nil {
		log.WithFields(log.Fields{
			"app": app,
		}).Fatal("Render function failed for app.")
	}

	//fmt.Println(overrideValues)

	cmd := []interface{}{
		"upgrade", app, chart,
		"--debug",
		"--dry-run",
		"--force",
		"--force",
		"--install",
		"--recreate-pods",
		"--values", "-",
	}

	//err := sh.Command("helm", cmd...).Run()
	err = sh.Command("helm", cmd...).SetInput(string(overrideValues)).Run()
	if err != nil {
		log.WithFields(log.Fields{
			"app":   app,
			"chart": chart,
		}).Fatal("Failed running `helm upgrade` for app.")
	}
}
