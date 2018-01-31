
# Kubernetes Survival Handbook

**Chapter 1: MultiHelm v0.4.0 at Minikube**

Contact: <josdotso@cisco.com>

## Develop a new MultiHelm App on your laptop, using Minikube

### Step 1: Install prerequisites.

[kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

[minikube](https://github.com/kubernetes/minikube)

[Helm](https://docs.helm.sh/using_helm/#quickstart)

[MultiHelm](https://***REMOVED***) (use tag `v0.4.0`)


### Step 2: Start Minikube

We'll use Minikube to run a small Kubernetes cluster on your laptop.

```
$ minikube start
Starting local Kubernetes v1.8.0 cluster...
Starting VM...
Getting VM IP address...
Moving files into cluster...
Setting up certs...
Connecting to cluster...
Setting up kubeconfig...
Starting cluster components...
Kubectl is now configured to use the cluster.
```

### Step 3: Confirm kubectl will interact with Minikube.

It may take a few minutes before Kubernetes is responsive.

```
$ kubectl config get-contexts
CURRENT  NAME               CLUSTER              AUTHINFO             NAMESPACE
*        minikube           minikube             minikube

$ kubectl config use-context minikube
Switched to context "minikube".

$ kubectl get nodes
NAME       STATUS    ROLES     AGE       VERSION
minikube   Ready     <none>    4m        v1.8.0

$ kubectl version
Client Version: version.Info{Major:"1", Minor:"8", GitVersion:"v1.8.2", GitCommit:"bdaeafa71f6c7c04636251031f93464384d54963", GitTreeState:"clean", BuildDate:"2017-10-24T21:07:53Z", GoVersion:"go1.9.1", Compiler:"gc", Platform:"darwin/amd64"}
Server Version: version.Info{Major:"1", Minor:"8", GitVersion:"v1.8.0", GitCommit:"0b9efaeb34a2fc51ff8e4d34ad9bc6375459c4a4", GitTreeState:"dirty", BuildDate:"2017-10-17T15:09:55Z", GoVersion:"go1.8.3", Compiler:"gc", Platform:"linux/amd64"}
```

### Step 4: Install Helm's Tiller service at Kubernetes.

```
$ helm init
$HELM_HOME has been configured at ${HOME}/.helm.

Tiller (the Helm server-side component) has been installed into your Kubernetes Cluster.
Happy Helming!
```

### Step 5: Confirm that Helm is ready to use with your Kubernetes.

```
$ helm version
Client: &version.Version{SemVer:"v2.7.0", GitCommit:"08c1144f5eb3e3b636d9775617287cc26e53dba4", GitTreeState:"clean"}
Server: &version.Version{SemVer:"v2.7.0", GitCommit:"08c1144f5eb3e3b636d9775617287cc26e53dba4", GitTreeState:"clean"}
```

### Step 6: Create a new git repository to hold your team's MultiHelm manifests

```
mkdir hello-multihelm
cd hello-multihelm
git init 
echo "# Hello MultiHelm" > README.md
mkdir -p apps configs
touch apps/.gitkeep configs/.gitkeep
git add .
git commit -a -m "First commit"
```

### Step 7: Take a look at an example chart's "values.yaml" file.

https://github.com/kubernetes/charts/tree/master/stable/wordpress

See how the defaults in
[stable/wordpress/values.yaml](https://github.com/kubernetes/charts/blob/master/stable/wordpress/values.yaml)
do not fit MiniKube? (example follows)

```
## Kubernetes configuration
## For minikube, set this to NodePort, elsewhere use LoadBalancer
##
serviceType: LoadBalancer
```

MultiHelm templates the **overriding** of Helm charts' `values.yaml` files.
In the next step, we'll use MultiHelm to override the `stable/wordpress`
chart's default values.

By the way, a "MultiHelm App" is basically a template for overriding
a Helm chart.

### Step 8: Create a MultiHelm App for "kubernetes/stable/wordpress"

Paste this at file `hello-multihelm/apps/wordpress.yaml` and be sure
to return to the `hello-multihelm` directory afterward.

The following is a MultHelm App template. We'll use a central
configuration file as the source of truth for it.

```
chart: stable/wordpress

# Version here means chart version. See `helm search wordpress`.
{{- if $app.version }}
version: {{ $app.version }}
{{- else }}
version: 0.8.2
{{- end }}

image: {{ $app.image }}:{{ $app.imageTag }}

imagePullPolicy: {{ .imagePullPolicy }}

serviceType: {{ $app.service.type }}
```

### Step 9: Create a MultiHelm configuration file for Minikube.

Paste this at file `hello-multihelm/configs/minikube.yaml` and be sure
to return to the `hello-multihelm` directory afterward.

The following is a MultHelm configuration file. We can centralize our
override variables for the `minkube` conext (or some other context) here.

```
targetContext: minikube

team: hello

maintainers:
  - your-email-here@example.com

apps:
  - name: wordpress
    alias: wordpress-blue
    key: .wordpressBlue

## For each app, the first appSource to find the app file exists as specified is selected.
## appSources are evaluated in the order declared here.
appSources:
  - name: apps
    kind: path
    source: ./apps
# - name: foo-deploy   # This is an example of how you might
#                      # implement shared MultiHelm apps via Git submodule.
#   kind: path
#   source: ../submodules/***REMOVED***/***REMOVED***/foo-deploy/multihelm/apps

imagePullPolicy: IfNotPresent

wordpressBlue:
  image: bitnami/wordpress
  imageTag: "4.9.0-r0"
  service:
    type: NodePort
```

### Step 10: Export the "MULTIHELM_CONFIG" environment variable..

```
$ export MULTIHELM_CONFIG="$(pwd)/configs/minikube.yaml"
$ echo ${MULTIHELM_CONFIG}
/tmp/hello-multihelm/configs/minikube.yaml
```

### Step 11: Print the rendered app (values override file) while simulating a MultiHelm install.

```
$ multihelm simulate --printRendered
```

### Step 12: Apply the HELM_CONFIG's apps (without verbose app rendering).

```
$ multihelm apply
```

### Step 13: Commit your changes to git.

```
git add .
git commit -a -m "Second commit."
```

### Step 14: Learn more at the MultiHelm readme.

https://***REMOVED***

Thanks for reading!

Feel free to email me with any comments, questions, PRs or requests.

I created a Spark channel for MultiHelm. Channel invites available upon request!

-Joshua M. Dotson <josdotso@cisco.com>
