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
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	log "github.com/sirupsen/logrus"
)

var (
	configFileFlag string
	versionNumber  string
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "multihelm",
	Short: "Operate multiple Helm charts",
	Long: `                   ___    __        __              ___
                  /\_ \  /\ \__  __/\ \            /\_ \
  ___ ___   __  __\//\ \ \ \ ._\/\_\ \ \___      __\//\ \     ___ ___
/. __. __.\/\ \/\ \ \ \ \ \ \ \/\/\ \ \  _  \  / __ \\ \ \  /  __. __.\
/\ \/\ \/\ \ \ \_\ \ \_\ \_\ \ \_\ \ \ \ \ \ \/\  __/ \_\ \_/\ \/\ \/\ \
\ \_\ \_\ \_\ \____/ /\____\\ \__\\ \_\ \_\ \_\ \____\/\____\ \_\ \_\ \_\
 \/_/\/_/\/_/\/___/  \/____/ \/__/ \/_/\/_/\/_/\/____/\/____/\/_/\/_/\/_/

MultiHelm simplifies multi-chart Helm workflows by rendering ephemeral Helm
chart override files based on templates populated with values from MultiHelm
YAML config files.

In other words: We heard you like templates, so we templated your Helm value
overrides.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Failed to execute RootCmd.")
	}
}

func init() {
	versionNumber = "v0.4.1"

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&configFileFlag, "config", "c", "",
		`config file (you can instead set MULTIHELM_CONFIG)`)
	RootCmd.PersistentFlags().BoolP("json", "j", false, "set logging to JSON format")

	// Beware that init() happens too early to read values from Viper...
	// See: https://github.com/spf13/cobra/issues/511
	//
	// TL;DR -- Use Viper for retrieving values, but read values no earlier than at "Run:" time.
	viper.BindPFlags(RootCmd.PersistentFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	var (
		configFile           string
		configFileEnv        string
		configFileEnvPresent bool
		configFileOrigin     string
	)

	// If environment variable is set, load its value.
	configFileEnv, configFileEnvPresent = os.LookupEnv("MULTIHELM_CONFIG")
	if configFileEnvPresent {
		configFile = configFileEnv
		configFileOrigin = "env"
	}
	if configFileFlag != "" {
		// If flag isn't empty, load its value instead.
		configFile = configFileFlag
		configFileOrigin = "flag"
	}
	viper.SetConfigFile(configFile)
	viper.SetEnvPrefix("multihelm") // will be uppercased automatically
	viper.AutomaticEnv()            // read in environment variables that match

	// If a configFile is found, read it in.
	err := viper.ReadInConfig()
	if viper.GetBool("json") {
		log.SetFormatter(&log.JSONFormatter{})
	}
	if err != nil {
		log.WithFields(log.Fields{
			"configFile":           configFile,
			"configFileEnv":        configFileEnv,
			"configFileEnvPresent": configFileEnvPresent,
			"configFileFlag":       configFileFlag,
			"configFileOrigin":     configFileOrigin,
			"configFileUsed":       getConfigFile(),
			"err":                  err,
		}).Warnln("Failed to load MultiHelm config.",
			"Please consider exporting environment variable: MULTIHELM_CONFIG.")
	}
}
