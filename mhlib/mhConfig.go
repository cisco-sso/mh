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

	"github.com/imdario/mergo"
)

// MHConfig is a set of options used during app deployment.
type MHConfig struct {
	Maintainers    []string `yaml:"maintainers"`
	PrintRendered  bool     `yaml:"printRendered"`
	NoRecreatePods bool     `yaml:"noRecreatePods"`
	Simulate       bool     `yaml:"simulate"`
	TargetContext  string   `yaml:"targetContext"`
	Team           string   `yaml:"team"`
	SETValues      []string
}

// DefaultMHConfig is the default mh config and will most likely be modified
// during execution with the following sources by rising priority:
//
// 1. This default configuration
// 2. The "mh" key in the configuration file
// 3. Environment variables starting with "MH_"
// 4. Command line flags
// 5. app-specific overrides in MH_CONFIG.
var DefaultMHConfig = MHConfig{
	Maintainers:    []string{"none"},
	PrintRendered:  false,
	NoRecreatePods: false,
	Simulate:       false,
	TargetContext:  "localhost",
	Team:           "sre",
	SETValues:      []string{""},
}

// MergeMHConfigs merges an arbitrary number of MHConfigs with rising priority.
func MergeMHConfigs(configs ...MHConfig) (*MHConfig, error) {
	if len(configs) < 2 {
		return nil, fmt.Errorf("Can't merge less than two configs")
	}

	result := MHConfig{}
	for _, config := range configs {
		if err := mergo.MergeWithOverwrite(&result, config); err != nil {
			return nil, err
		}
	}

	return &result, nil
}
