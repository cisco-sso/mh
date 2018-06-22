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
	"path"
	"path/filepath"
	"strings"
)

// AppSourceConfig is what can be defined in a mh configuration file and is used
// create an AppSource struct. It defines where mh looks for app files.
//
// Todo: Define kinds of sources here.
type AppSourceConfig struct {
	Kind   string `yaml:"kind"`
	Name   string `yaml:"name"`
	Source string `yaml:"source"`
}

// AppSourceConfigs is an Array of AppSourceConfigs defined in a mh configuration file
type AppSourceConfigs []AppSourceConfig

// AppSource contains an AppSourceConfig and all AppFiles retrieved via it.
type AppSource struct {
	AppSourceConfig
	Files AppFiles
}

// File takes an AppConfig and sets its File to the matching entry of the
// AppSource if one exists. Returns false if none exists.
func (as *AppSource) File(appConfig *AppConfig) bool {
	if file, ok := as.Files[appConfig.Name]; ok {
		appConfig.File = &file

		return true
	}

	return false
}

// NewAppSource returns an AppSource based on a given AppSourceConfig and path to
// mh configuration(for "configPath" kind of AppSources).
//
// Todo: Get rid of "configFile" argument in favor of defined kinds of
// AppSources - see AppSourceConfig.
func NewAppSource(config AppSourceConfig, configFile string) (*AppSource, error) {
	// Build glob pattern depending on kind of AppSource
	var pattern string
	if config.Kind == "path" {
		pattern = config.Source + "/*.y*ml"
	} else if config.Kind == "configPath" {
		pattern = path.Dir(configFile) + config.Source + "/*.y*ml"
	}

	// Get all files the glob pattern matches
	filePaths, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("Failed to determine file from AppSource: %s", err.Error())
	}

	// Create an AppFile for each matched file
	files := AppFiles{}
	for _, path := range filePaths {
		base := filepath.Base(path)
		name := strings.TrimSuffix(base, filepath.Ext(base))
		localpath := path

		file := AppFile{
			Path: &localpath,
		}
		files[name] = file
	}

	return &AppSource{
		config,
		files,
	}, nil
}
