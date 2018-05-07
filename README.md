
# mh

```
$ mh --help
            __
           /\ \
  ___ ___  \ \ \___
/. __. __.\ \ \  _  \
/\ \/\ \/\ \ \ \ \ \ \
\ \_\ \_\ \_\ \ \_\ \_\
 \/_/\/_/\/_/  \/_/\/_/

mh simplifies multi-chart Helm workflows by rendering ephemeral Helm
chart override files based on templates populated with values from mh
YAML config files.

In other words: We heard you like templates, so we templated your Helm value
overrides.

Usage:
  mh [command]

Available Commands:
  apply       Apply apps
  destroy     Destroy apps
  help        Help about any command
  license     Print license information.
  simulate    Simulate apps
  status      Get status of apps
  version     Print version information.

Flags:
  -c, --config string     config file (you can instead set MH_CONFIG)
  -h, --help              help for mh
  -j, --json              set logging to JSON format

Use "mh [command] --help" for more information about a command.
```

```
$ mh license
Copyright © 2018 Cisco Systems, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

## Getting Started

### Install mh.

(NOTE: Build below currently requires a working golang dev environment.)

(Simplified install coming soon!)

```
go get -u github.com/golang/dep/cmd/dep   # https://github.com/golang/dep
mkdir -p $GOPATH/src/github.com/cisco-sso
cd $GOPATH/src/github.com/cisco-sso
git clone git@github.com:cisco-sso/mh.git
cd mh
dep ensure
go build -o /tmp/mh main.go
sudo cp /tmp/mh /usr/local/bin
rm -f /tmp/mh
```

### Select a kubectl context.

(In general, this is when you "choose a Kubernetes cluster" to manage.)

```
kubectl config get-contexts
kubectl config use-context minikube
```

### Select a mh config.

(There's usally one mh config per kubetl context, but we've left it open
ended so that multple teams can more easily work together on one cluster.)

```
export MH_CONFIG="/path/to/clusters/minikube/mh/main.yaml"
```

### Get status of everything at context "minkube" managed by this mh config.

(This basically runs `helm status` for each app you target.)

```
mh status
# ^ get status for all apps in `main.yaml`

mh status wordpress
# ^ get status for just these app(s)
```

### Simulate app upgrades (or simulate install of apps, as needed).

(For each app you target, simulate runs a Helm upgrade/install
with debug and dry-run modes enabled.)

```
mh simulate
# ^ simulate install/upgrade for all apps in `main.yaml`

mh simulate --printRendered
# ^ simulate install/upgrade for all apps in `main.yaml`
#   (verbosely printing app template renderings)

mh simulate wordpress
# ^ simulate install/upgrade just these app(s),
#   even if they are not in `main.yaml`
```

### Apply app upgrades (or install apps, as needed).

(For each app you target, apply runs a Helm upgrade/install).

```
mh apply
# ^ apply install/upgrade for all apps in `main.yaml`

mh apply --printRendered
# ^ apply install/upgrade for all apps in `main.yaml`
#   (verbosely printing app template renderings)

mh apply wordpress
# ^ apply install/upgrade just these app(s),
#   even if they are not in `main.yaml`
```

### Destroy apps (if they are known to Helm).

(For each app you target, apply runs a Helm delete without purge).

```
mh destroy
# ^ destroy all apps in `main.yaml`

mh destroy wordpress
# ^ destroy just these app(s),
#   even if they are not in `main.yaml`
```

### Log to JSON!

```
mh status foo --json 2>&1 | jq --slurp
```

## Contributing

For information on contributing to this repo, please see the accompanying `CONTRIBUTING.md` file.
