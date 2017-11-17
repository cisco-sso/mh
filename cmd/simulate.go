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
	"fmt"

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
		lateInit("simulate")
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

	simulateCmd.PersistentFlags().BoolP("printRendered", "p", false, "print rendered override values")
	viper.BindPFlags(simulateCmd.PersistentFlags())
}

func simulate(app string) {

	chart, overrideValues, err := render(app)
	if err != nil {
		log.WithFields(log.Fields{
			"app": app,
		}).Fatal("Render function failed for app.")
	}

	if viper.GetBool("printRendered") {
		log.WithFields(log.Fields{
			"app":   app,
			"chart": chart,
		}).Info("Printing rendered override values for app.")
		fmt.Print(string(overrideValues))
		log.WithFields(log.Fields{
			"app":   app,
			"chart": chart,
		}).Info("Done printing rendered override values for app.")
	}

	cmd := []interface{}{
		"upgrade", app, chart,
		"--debug",
		"--dry-run",
		"--force",
		"--install",
		"--recreate-pods",
		"--values", "-",
	}

	err = sh.Command("helm", cmd...).SetInput(string(overrideValues)).Run()
	if err != nil {
		log.WithFields(log.Fields{
			"app":   app,
			"chart": chart,
		}).Fatal("Failed running (simulated) `helm upgrade` for app.")
	}
}
