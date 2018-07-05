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

package mhlib

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

// MHConfigFile is the structure of a mh configuration file.
type MHConfigFile struct {
	MH         MHConfig         `yaml:"mh"`
	Apps       AppConfigs       `yaml:"apps"`
	AppSources AppSourceConfigs `yaml:"appSources"`
}

// EffectiveApps returns all Apps that are configured in a MHConfigFile,
// optionally filtering them with given expressions. It matches them against the
// MHConfigFiles AppSourceConfigs, if their File is not overridden. Also passes
// a given logger and effective MHConfig down to them.
//
// Todo: Support more than simple names as filter. Improve the below algorithm.
func (c *MHConfigFile) EffectiveApps(logger *logrus.Entry, configFile string, filters []string, effectiveMHConfig MHConfig) (*Apps, error) {
	var effectiveAppSources AppSources
	var effectiveApps Apps

	// Build effective app sources from AppSourceConfigs defined in configuration
	// file.
	for _, appSourceConfig := range c.AppSources {
		appSource, err := NewAppSource(appSourceConfig, configFile)
		if err != nil {
			return nil, err
		}

		if len(appSource.Files) == 0 {
			logger.WithField("appSource", appSource.Name).Warn("AppSource matches no files")
		}

		effectiveAppSources = append(effectiveAppSources, *appSource)
	}

	// Build effective apps from AppConfigs given in configuration file, matching
	// them with effective AppSources and MHConfig built above.
	for _, appConfig := range c.Apps {
		// If no filters are defined, add the app immediately
		if len(filters) == 0 {
			// Match the app config with configured app sources if File is not
			// overridden.
			if appConfig.File == nil {
				if !effectiveAppSources.File(&appConfig) {
					logger.WithFields(logrus.Fields{
						"app": appConfig.Name,
					}).Error("App not found in sources")
					return nil, fmt.Errorf("App not found in sources")
				}
			}

			app, err := NewApp(logger, appConfig, effectiveMHConfig)
			if err != nil {
				return nil, err
			}

			effectiveApps = append(effectiveApps, *app)
		} else {
			// If filters are defined, try to match against them
			for _, filter := range filters {
				// Add the app only if the filter matches its Name or Alias
				if appConfig.Name == filter || appConfig.Alias == filter {
					logger.WithFields(logrus.Fields{
						"app":    appConfig.Name,
						"filter": filter,
					}).Info("App matched filter")

					// Match the app config with configured app sources if File is not
					// overridden.
					if appConfig.File == nil {
						if !effectiveAppSources.File(&appConfig) {
							logger.WithFields(logrus.Fields{
								"app": appConfig.Name,
							}).Error("App not found in sources")
							return nil, fmt.Errorf("App not found in sources")
						}
					}

					app, err := NewApp(logger, appConfig, effectiveMHConfig)
					if err != nil {
						logger.WithFields(logrus.Fields{
							"app":   appConfig.Name,
							"error": err,
						}).Error("App creation failed")
						return nil, err
					}

					effectiveApps = append(effectiveApps, *app)
					break
				}
			}
		}
	}

	return &effectiveApps, nil
}
