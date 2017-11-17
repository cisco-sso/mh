
# MultiHelm

```
$ multihelm --help
                   ___    __        __              ___
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
overrides.

Usage:
  multihelm [command]

Available Commands:
  apply       Apply apps
  destroy     Destroy apps
  help        Help about any command
  license     Print license information.
  simulate    Simulate apps
  status      Get status of apps
  version     Print version information.

Flags:
  -a, --appsPath string   apps path (default "./apps")
  -c, --config string     config file (you can instead set MULTIHELM_CONFIG)
  -h, --help              help for multihelm

Use "multihelm [command] --help" for more information about a command.
```

```
$ multihelm license
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
```

## Kubernetes Survival Handbook

This guide to Kubernetes lives in the `docs` folder of this repo.

Chapter 1, entitled "MultiHelm at Minikube", starts [here](https://***REMOVED***/browse/docs/KubernetesSurvivalHandbook/chapter1.md).


## Getting Started

### Install MultiHelm.

(NOTE: Build below currently requires a working golang dev environment.)

(Simplified install coming soon!)

```
git clone ssh://***REMOVED***/***REMOVED***/multihelm.git
cd multihelm
go build main.go -o /tmp/multihelm
sudo cp /tmp/multihelm /usr/local/bin
rm -f /tmp/multihelm
```

### Clone the example "hello-multihelm" repo.

(Or make a new MultiHelm manifests repo for your team.)

https://***REMOVED***/bitbucket/projects/***REMOVED***/repos/hello-multihelm

```
git clone ssh://***REMOVED***/***REMOVED***/hello-multihelm.git
cd hello-multihelm
```

### Initialize and update git submodules used by this repo.

(It is a best practice to keep charts external to your MultHelm and version-lock
your usage of them via git submodule.)

```
git submodule init
git submodule update
```

### Select a kubectl context.

(In general, this is when you "choose a Kubernetes cluster" to manage.)

```
kubectl config get-contexts
kubectl config use-context minikube
```

### Select a MultiHelm config.

(There's usally one MultiHelm config per kubetl context, but we've left it open
ended so that multple teams can more easily work together on one cluster.)

```
export MULTIHELM_CONFIG="$(pwd)/configs/minikube.yaml"
```

### Get status of everything at context "minkube" managed by this MultiHelm config.

(This basically runs `helm status` for each app you target.)

```
multihelm status
# ^ get status for all apps in `minikube.yaml`

multihelm status wordpress
# ^ get status for just these app(s)
```

### Simulate app upgrades (or simulate install of apps, as needed).

(For each app you target, simulate runs a Helm upgrade/install
with debug and dry-run modes enabled.)

```
multihelm simulate
# ^ simulate install/upgrade for all apps in `minikube.yaml`

multihelm simulate --printRendered
# ^ simulate install/upgrade for all apps in `minikube.yaml`
#   (verbosely printing app template renderings)

multihelm simulate wordpress
# ^ simulate install/upgrade just these app(s),
#   even if they are not in `minikube.yaml`
```

### Apply app upgrades (or install apps, as needed).

(For each app you target, apply runs a Helm upgrade/install).

```
multihelm apply
# ^ apply install/upgrade for all apps in `minikube.yaml`

multihelm apply --printRendered
# ^ apply install/upgrade for all apps in `minikube.yaml`
#   (verbosely printing app template renderings)

multihelm apply wordpress
# ^ apply install/upgrade just these app(s),
#   even if they are not in `minikube.yaml`
```

### Destroy apps (if they are known to Helm).

(For each app you target, apply runs a Helm delete without purge).

```
multihelm destroy
# ^ destroy all apps in `minikube.yaml`

multihelm destroy wordpress
# ^ destroy just these app(s),
#   even if they are not in `minikube.yaml`
```
