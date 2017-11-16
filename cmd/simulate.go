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
	//"github.com/codeskyblue/go-sh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// simulateCmd represents the simulate command
var simulateCmd = &cobra.Command{
	Use:   "simulate",
	Short: "Simulate apps",
	Long: `Simulate the apply of one or more MultiHelm apps. If you do not specify one or more
apps, MultiHelm acts on all apps in your MultiHelm config.`,
	Run: func(cmd *cobra.Command, args []string) {
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

	renderedValues, err := render(app)
	if err != nil {
		os.Exit(1)
	}

	fmt.Println(renderedValues)
	/*
		// Render override file
		cfgFile := viper.ConfigFileUsed()

		appsPath := "./apps"
		chartPath := "./src"
		cmd := []interface{}{
			"upgrade", "--install", "--force", "--recreate-pods",
			"--debug", "--dry-run", "--force", app, chartPath,
		}
		err := sh.Command("helm", cmd...).Run()
		if err != nil {
			os.Exit(1)
		}
	*/
}
