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

import log "github.com/sirupsen/logrus"

type AppLog struct {
	app           *App
	appFile       string
	appSources    []AppSource
	cmd           []interface{}
	configFile    string
	data          []byte
	err           error
	id            string
	method        string
	printRendered bool
	purge         bool
	reason        string
	simulate      bool
}

func (a *AppLog) Error() {
	var id string
	if a.id == "" {
		id = a.app.Id()
	} else {
		id = a.id
	}
	msg := "Failed running '" + a.method + "' for app '" + id + "'"
	log.WithFields(log.Fields{
		"app":        a.app,
		"appFile":    a.appFile,
		"appSources": a.appSources,
		"cmd":        a.cmd,
		"data":       string(a.data),
		"err":        a.err,
		"id":         id,
		"method":     a.method,
		"purge":      a.purge,
		"reason":     a.reason,
		"simulate":   a.simulate,
	}).Fatal(msg)
}

func (a *AppLog) Info(msg string) {
	var id string
	if a.id == "" {
		id = a.app.Id()
	} else {
		id = a.id
	}
	log.WithFields(log.Fields{
		"app":        a.app,
		"appFile":    a.appFile,
		"appSources": a.appSources,
		"cmd":        a.cmd,
		"err":        a.err,
		"data":       string(a.data),
		"id":         id,
		"method":     a.method,
		"purge":      a.purge,
		"reason":     a.reason,
		"simulate":   a.simulate,
	}).Info(msg)
}
