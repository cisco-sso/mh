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
	cfgFile          string
	cfgFileSetMethod string
	currentContext   string
	tryCfgFile       string
	versionNumber    string
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
	versionNumber = "v0.1.1"
	currentContext = getCurrentContext()

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", `config file (you can instead set MULTIHELM_CONFIG)`)
	RootCmd.PersistentFlags().StringP("appsPath", "a", "./apps", "apps path")

	// Beware that init() happens too early to read values from Viper...
	// See: https://github.com/spf13/cobra/issues/511
	//
	// TL;DR -- Use Viper for retrieving values, but read values no earlier than at "Run:" time.
	viper.BindPFlags(RootCmd.PersistentFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// If env var is set, update cfgFile
	envCfgFile, present := os.LookupEnv("MULTIHELM_CONFIG")
	if present {
		tryCfgFile = envCfgFile
		cfgFileSetMethod = "env"
	}
	if cfgFile != "" {
		// Override config file from the flag.
		tryCfgFile = cfgFile
		cfgFileSetMethod = "flag"
	}
	viper.SetConfigFile(tryCfgFile)
	viper.SetEnvPrefix("multihelm") // will be uppercased automatically
	viper.AutomaticEnv()            // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"err":              err,
			"tryCfgFile":       tryCfgFile,
			"cfgFileSetMethod": cfgFileSetMethod,
		}).Fatalln("Failed to load MultiHelm config.",
			"Please consider exporting env var MULTIHELM_CONFIG.")

	}
}
