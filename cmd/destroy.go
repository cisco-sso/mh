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
	lib "github.com/cisco-sso/mh/mhlib"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// destroyCmd represents the destroy command
var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy apps",
	Long: `Destroy one or more mh apps. If you do not specify one or more
apps, mh acts on all apps in your mh config.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := logrus.New().WithField("command", "destroy")
		if viper.GetBool("json") {
			logger.Logger.Formatter = new(logrus.JSONFormatter)
		}
		mhConfigFile := unmarshalConfig(logger)

		// Merge configuration from file, environment and CLI into default
		// configuration
		effectiveMHConfig, err := lib.MergeMHConfigs(lib.DefaultMHConfig, mhConfigFile.MH)
		if err != nil {
			logger.WithField("error", err).Fatal("Failed to build effective MH configuration")
		}

		// Ensure TargetContext is the current kubectl context
		ensureCurrentContext(logger, *effectiveMHConfig)

		// Get effective apps
		apps, err := mhConfigFile.EffectiveApps(logger, viper.ConfigFileUsed(), args, *effectiveMHConfig)
		if err != nil {
			logger.WithField("error", err).Fatal("Failed to build effective apps")
		}

		if err := apps.Destroy(); err != nil {
			logger.Fatal("Failed running destroy")
		}
	},
}

func init() {
	RootCmd.AddCommand(destroyCmd)

	destroyCmd.PersistentFlags().BoolP("purge", "p", false, "purge this app from Helm Tiller")
	viper.BindPFlags(destroyCmd.PersistentFlags())
}
