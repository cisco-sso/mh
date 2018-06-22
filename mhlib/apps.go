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
	"github.com/sirupsen/logrus"
)

// Apps is an array of apps.
type Apps []App

// Apply runs Apply on each App
func (a Apps) Apply(configFile string) error {
	for _, app := range a {
		cmd, err := app.Apply(configFile)
		if err != nil {
			app.log.WithFields(logrus.Fields{
				"app":   app.Name,
				"cmd":   cmd,
				"error": err,
			}).Fatal("Failed running apply")

			return err
		}

	}

	return nil
}

// Destroy runs Destroy on each App
func (a Apps) Destroy() error {
	for _, app := range a {
		cmd, err := app.Destroy(false)
		if err != nil {
			app.log.WithFields(logrus.Fields{
				"app":   app.Name,
				"cmd":   cmd,
				"error": err,
			}).Fatal("Failed running destroy")

			return err
		}
	}

	return nil
}

// Simulate runs Simulate on each App
func (a Apps) Simulate(configFile string) error {
	for _, app := range a {
		cmd, err := app.Simulate(configFile)
		if err != nil {
			app.log.WithFields(logrus.Fields{
				"app":   app.Name,
				"cmd":   cmd,
				"error": err,
			}).Fatal("Failed running simulate")

			return err
		}
	}

	return nil
}
