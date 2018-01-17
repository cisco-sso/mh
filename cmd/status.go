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

import "github.com/spf13/cobra"

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status [APP]...",
	Short: "Get status of apps",
	Long: `Get status one or more MultiHelm apps. If you do not specify one or more
apps, MultiHelm acts on all apps in your MultiHelm config.`,
	Run: func(cmd *cobra.Command, args []string) {
		lateInit("status")

		apps := getApps(args)

		apps.Status()
	},
}

func init() {
	RootCmd.AddCommand(statusCmd)
}
