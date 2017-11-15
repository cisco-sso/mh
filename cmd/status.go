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
	"os"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status [APP]...",
	Short: "Get status of apps",
	Long: `Get status one or more MultiHelm apps. If you do not specify one or more
apps, MultiHelm will get status for all apps in your context config.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			for _, arg := range args {
				status(arg)
			}
		} else {
			for _, arg := range viper.GetStringSlice("apps") {
				status(arg)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// statusCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// statusCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func status(app string) {
	err := sh.Command("helm", "status", app).Run()
	if err != nil {
		os.Exit(1)
	}
}
