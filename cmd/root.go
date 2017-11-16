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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
)

var cfgFile string
var currentContext string
var versionNumber string

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "multihelm",
	Short: "Operate multiple Helm charts",
	Long: `MultiHelm simplifies multi-chart workflows by rendering ephemeral Helm
chart override files based on templates populated with values from
MultiHelm YAML config files.`,
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
	versionNumber = "v0.1.0"
	currentContext = getCurrentContext()

	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", `config file (default "${HOME}/.multihelm.yaml")`)
	RootCmd.PersistentFlags().StringP("appsPath", "", "./apps", "apps path")
	RootCmd.PersistentFlags().StringP("targetContext", "", "minikube", "target kubectl context")

	// Beware that init() happens too early to read values from Viper...
	// See: https://github.com/spf13/cobra/issues/511
	//
	// TL;DR -- Use Viper for retrieving values, but read values no earlier than at "Run:" time.
	viper.BindPFlags(RootCmd.PersistentFlags())
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			log.WithFields(log.Fields{
				"err": err,
			}).Fatal("Failed to lookup home directory.")
		}

		// Search config in home directory with name ".multihelm" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".multihelm")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Failed to load MultiHelm config.")

	}
}
