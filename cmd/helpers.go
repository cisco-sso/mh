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
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/codeskyblue/go-sh"
	"github.com/spf13/viper"

	lib "github.com/cisco-sso/mh/mhlib"
)

// ensureCurrentContext compares the current kubectl context with the targent
// context and exits if they do not match.
func ensureCurrentContext(logger *logrus.Entry, config lib.MHConfig) {
	// Fetch current context
	out, err := sh.Command("kubectl", "config", "current-context").Output()
	if err != nil {
		logger.WithField("error", err).Fatal("Failed running `kubectl config current-context`.")
	}
	currentContext := strings.TrimSuffix(string(out), "\n")

	// Compare contexts
	if config.TargetContext != currentContext {
		logger.WithFields(logrus.Fields{
			"currentContext": currentContext,
			"targetContext":  config.TargetContext,
		}).Fatal("`kubectl config current-context` does not match configured targetContext")
	}
}

// unmarshalConfig creates a MHConfigFile struct from MH_CONFIG
func unmarshalConfig(logger *logrus.Entry) lib.MHConfigFile {
	mhConfigFile := lib.MHConfigFile{}
	if err := viper.Unmarshal(&mhConfigFile); err != nil {
		logger.WithField("error", err).Fatal("Failed to unmarshal mh configuration file")
	}
	logger = logger.WithField("configFile", viper.ConfigFileUsed())

	return mhConfigFile
}
