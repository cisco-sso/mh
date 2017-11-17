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

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status [APP]...",
	Short: "Get status of apps",
	Long: `Get status one or more MultiHelm apps. If you do not specify one or more
apps, MultiHelm acts on all apps in your MultiHelm config.`,
	Run: func(cmd *cobra.Command, args []string) {
		lateInit("status")
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
}

func status(app string) {
	err := sh.Command("helm", "status", app).Run()
	if err != nil {
		log.WithFields(log.Fields{
			"app": app,
			"err": err,
		}).Fatal("Failed to get status for app.")
	}
}
