
# Kubernetes Survival Handbook

**Chapter 1: MultiHelm at Minikube**

Contact: <josdotso@cisco.com>

## Develop a new MultiHelm App on your laptop, using Minikube

### Step 1: Install prerequisites.

[kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/)

[minikube](https://github.com/kubernetes/minikube)

[Helm](https://docs.helm.sh/using_helm/#quickstart)

[MultiHelm](https://***REMOVED***)


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
mkdir -p apps configs src
touch apps/.gitkeep configs/.gitkeep src/.gitkeep
git add .
git commit -a -m "First commit"
```

### Step 7: Add a "kubernetes/charts" repo as a submodule.

We recommend that you use a fork of [kubernetes/charts](https://github.com/kubernetes/charts)
for this submodule. We're using the unforked repo only for this example.

I am happy to help you create such a fork if you'd like. -josdosto

Note that we namespace our submodules in a Golang-style tree. This allows
you to use two different forks of the same repo at the same time. It also
helps to reduce confusion and submodule name collisions.

```
git submodule add -b master ssh://git@github.com:kubernetes/charts.git src/github.com/kubernetes/charts
git add .
git commit -a -m "Add submodule 'src/github.com/kubernetes/charts'"
```

### Step 8: Take a look at an example chart's "values.yaml" file.

We'll use MultiHelm to override these values. MultiHelm templates the **overriding** of
Helm charts' `values.yaml` files.

It's also wise to review the `README.md` for the chart you'd like to make into a
"MultiHelm App".

"MultiHelm App" is more accurately described as "A template for overriding a
particular Helm chart".

```
$ ls src/github.com/kubernetes/charts/stable/
Display all 109 possibilities? (y or n)
acs-engine-autoscaler/  hadoop/                 moodle/                 selenium/
artifactory/            heapster/               msoms/                  sensu/
aws-cluster-autoscaler/ influxdb/               mysql/                  sentry/
buildkite/              ipfs/                   namerd/                 sonarqube/
centrifugo/             jasperreports/          nginx-ingress/          sonatype-nexus/
chaoskube/              jenkins/                nginx-lego/             spark/
chronograf/             joomla/                 odoo/                   spartakus/
cluster-autoscaler/     kapacitor/              opencart/               spinnaker/
cockroachdb/            keel/                   openvpn/                spotify-docker-gc/
concourse/              kube-lego/              orangehrm/              stash/
consul/                 kube-ops-view/          osclass/                sugarcrm/
coredns/                kube-state-metrics/     owncloud/               suitecrm/
coscale/                kube2iam/               parse/                  sumokube/
dask-distributed/       kubed/                  percona/                sumologic-fluentd/
datadog/                kubernetes-dashboard/   phabricator/            swift/
dokuwiki/               linkerd/                phpbb/                  sysdig/
drupal/                 locust/                 postgresql/             telegraf/
etcd-operator/          magento/                prestashop/             testlink/
external-dns/           mailhog/                prometheus/             traefik/
factorio/               mariadb/                rabbitmq/               uchiwa/
fluent-bit/             mcrouter/               redis/                  voyager/
g2/                     mediawiki/              redis-ha/               weave-cloud/
gcloud-endpoints/       memcached/              redmine/                wordpress/
gcloud-sqlproxy/        metabase/               rethinkdb/              zeppelin/
ghost/                  minecraft/              risk-advisor/           zetcd/
gitlab-ce/              minio/                  rocketchat/
gitlab-ee/              mongodb/                sapho/
grafana/                mongodb-replicaset/     searchlight/ 

$ ls src/github.com/kubernetes/charts/stable/wordpress
Chart.yaml  README.md  requirements.lock  requirements.yaml  templates  values.yaml

$ cat src/github.com/kubernetes/charts/stable/wordpress/values.yaml
```

(`values.yaml` follows)

```
## Bitnami WordPress image version
## ref: https://hub.docker.com/r/bitnami/wordpress/tags/
##
image: bitnami/wordpress:4.9.0-r0

## Specify a imagePullPolicy
## ref: http://kubernetes.io/docs/user-guide/images/#pre-pulling-images
##
imagePullPolicy: IfNotPresent

## User of the application
## ref: https://github.com/bitnami/bitnami-docker-wordpress#environment-variables
##
wordpressUsername: user

## Application password
## Defaults to a random 10-character alphanumeric string if not set
## ref: https://github.com/bitnami/bitnami-docker-wordpress#environment-variables
##
# wordpressPassword:

## Admin email
## ref: https://github.com/bitnami/bitnami-docker-wordpress#environment-variables
##
wordpressEmail: user@example.com

## First name
## ref: https://github.com/bitnami/bitnami-docker-wordpress#environment-variables
##
wordpressFirstName: FirstName

## Last name
## ref: https://github.com/bitnami/bitnami-docker-wordpress#environment-variables
##
wordpressLastName: LastName

## Blog name
## ref: https://github.com/bitnami/bitnami-docker-wordpress#environment-variables
##
wordpressBlogName: User's Blog!

## Set to `yes` to allow the container to be started with blank passwords
## ref: https://github.com/bitnami/bitnami-docker-wordpress#environment-variables
allowEmptyPassword: yes

## SMTP mail delivery configuration
## ref: https://github.com/bitnami/bitnami-docker-wordpress/#smtp-configuration
##
# smtpHost:
# smtpPort:
# smtpUser:
# smtpPassword:
# smtpUsername:
# smtpProtocol:

##
## MariaDB chart configuration
##
mariadb:
  ## MariaDB admin password
  ## ref: https://github.com/bitnami/bitnami-docker-mariadb/blob/master/README.md#setting-the-root-password-on-first-run
  ##
  # mariadbRootPassword:

  ## Create a database
  ## ref: https://github.com/bitnami/bitnami-docker-mariadb/blob/master/README.md#creating-a-database-on-first-run
  ##
  mariadbDatabase: bitnami_wordpress

  ## Create a database user
  ## ref: https://github.com/bitnami/bitnami-docker-mariadb/blob/master/README.md#creating-a-database-user-on-first-run
  ##
  mariadbUser: bn_wordpress

  ## Password for mariadbUser
  ## ref: https://github.com/bitnami/bitnami-docker-mariadb/blob/master/README.md#creating-a-database-user-on-first-run
  ##
  # mariadbPassword:

  ## Enable persistence using Persistent Volume Claims
  ## ref: http://kubernetes.io/docs/user-guide/persistent-volumes/
  ##
  persistence:
    enabled: true
    ## mariadb data Persistent Volume Storage Class
    ## If defined, storageClassName: <storageClass>
    ## If set to "-", storageClassName: "", which disables dynamic provisioning
    ## If undefined (the default) or set to null, no storageClassName spec is
    ##   set, choosing the default provisioner.  (gp2 on AWS, standard on
    ##   GKE, AWS & OpenStack)
    ##
    # storageClass: "-"
    accessMode: ReadWriteOnce
    size: 8Gi

## Kubernetes configuration
## For minikube, set this to NodePort, elsewhere use LoadBalancer
##
serviceType: LoadBalancer

## Allow health checks to be pointed at the https port
healthcheckHttps: false

## Configure ingress resource that allow you to access the
## Wordpress instalation. Set up the URL
## ref: http://kubernetes.io/docs/user-guide/ingress/
##
ingress:
  ## Set to true to enable ingress record generation
  enabled: false

  ## The list of hostnames to be covered with this ingress record.
  ## Most likely this will be just one host, but in the event more hosts are needed, this is an array
  hosts:
  - name: wordpress.local

    ## Set this to true in order to enable TLS on the ingress record
    ## A side effect of this will be that the backend wordpress service will be connected at port 443
    tls: false

    ## If TLS is set to true, you must declare what secret will store the key/certificate for TLS
    tlsSecret: wordpress.local-tls

    ## Ingress annotations done as key:value pairs
    ## If you're using kube-lego, you will want to add:
    ## kubernetes.io/tls-acme: true
    ##
    ## For a full list of possible ingress annotations, please see
    ## ref: https://github.com/kubernetes/ingress-nginx/blob/master/docs/annotations.md
    ##
    ## If tls is set to true, annotation ingress.kubernetes.io/secure-backends: "true" will automatically be set
    annotations:
    #  kubernetes.io/ingress.class: nginx
    #  kubernetes.io/tls-acme: true

  secrets:
  ## If you're providing your own certificates, please use this to add the certificates as secrets
  ## key and certificate should start with -----BEGIN CERTIFICATE----- or
  ## -----BEGIN RSA PRIVATE KEY-----
  ##
  ## name should line up with a tlsSecret set further up
  ## If you're using kube-lego, this is unneeded, as it will create the secret for you if it is not set
  ##
  ## It is also possible to create and manage the certificates outside of this helm chart
  ## Please see README.md for more information
  # - name: wordpress.local-tls
  #   key:
  #   certificate:

## Enable persistence using Persistent Volume Claims
## ref: http://kubernetes.io/docs/user-guide/persistent-volumes/
##
persistence:
  enabled: true
  ## wordpress data Persistent Volume Storage Class
  ## If defined, storageClassName: <storageClass>
  ## If set to "-", storageClassName: "", which disables dynamic provisioning
  ## If undefined (the default) or set to null, no storageClassName spec is
  ##   set, choosing the default provisioner.  (gp2 on AWS, standard on
  ##   GKE, AWS & OpenStack)
  ##
  # storageClass: "-"
  accessMode: ReadWriteOnce
  size: 10Gi

## Configure resource requests and limits
## ref: http://kubernetes.io/docs/user-guide/compute-resources/
##
resources:
  requests:
    memory: 512Mi
    cpu: 300m

## Node labels for pod assignment
## Ref: https://kubernetes.io/docs/user-guide/node-selection/
##
nodeSelector: {}
```

(Note how the defaults don't favor Minikube.)

excerpt (from the file above):
```
## Kubernetes configuration
## For minikube, set this to NodePort, elsewhere use LoadBalancer
##
serviceType: LoadBalancer
```

**Let's use MultiHelm to override these defaults!**

### Step 8: Create a MultiHelm App for "kubernetes/stable/wordpress"

Paste this at file `hello-multihelm/apps/wordpress.yaml` and be sure
to return to the `hello-multihelm` directory afterward.

The following is a MultHelm App template. We'll use a central
configuration file as the source of truth for it.

```
---
chart: ./src/github.com/kubernetes/charts/stable/wordpress

{{- $name := "wordpress" }}
{{- $app := .wordpress }}

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
---
targetContext: minikube

appsPath: ./apps

team: hello

maintainers:
  - your-email-here@example.com

apps:
  - wordpress

imagePullPolicy: IfNotPresent

wordpress:
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
INFO[0000] Initializing MultiHelm.                       appsPath=./apps cfgFile=/tmp/hello-multihelm/configs/minikube.yaml currentContext=minikube targetContext=minikube versionNumber=v0.1.0
INFO[0000] `simulate` called.
INFO[0000] Rendering app.                                app=wordpress
INFO[0000] Building chart dependencies for app.          app=wordpress chart=./src/github.com/kubernetes/charts/stable/wordpress err="<nil>"
INFO[0001] Printing rendered override values for app.    app=wordpress chart=./src/github.com/kubernetes/charts/stable/wordpress
---
chart: ./src/github.com/kubernetes/charts/stable/wordpress

image: bitnami/wordpress:4.9.0-r0

imagePullPolicy: IfNotPresent

serviceType: NodePort
INFO[0001] Done printing rendered override values for app.  app=wordpress chart=./src/github.com/kubernetes/charts/stable/wordpress
[debug] Created tunnel using local port: '61268'

[debug] SERVER: "localhost:61268"

Release "wordpress" does not exist. Installing it now.
[debug] CHART PATH: /tmp/hello-multihelm/src/github.com/kubernetes/charts/stable/wordpress

NAME:   wordpress
REVISION: 1
RELEASED: Thu Nov 16 20:43:56 2017
CHART: wordpress-0.7.2
USER-SUPPLIED VALUES:
chart: ./src/github.com/kubernetes/charts/stable/wordpress
image: bitnami/wordpress:4.9.0-r0
imagePullPolicy: IfNotPresent
serviceType: NodePort

COMPUTED VALUES:
allowEmptyPassword: true
chart: ./src/github.com/kubernetes/charts/stable/wordpress
healthcheckHttps: false
image: bitnami/wordpress:4.9.0-r0
imagePullPolicy: IfNotPresent
ingress:
  enabled: false
  hosts:
  - annotations: null
    name: wordpress.local
    tls: false
    tlsSecret: wordpress.local-tls
  secrets: null
mariadb:
  global: {}
  image: bitnami/mariadb:10.1.23-r2
  imagePullPolicy: IfNotPresent
  mariadbDatabase: bitnami_wordpress
  mariadbUser: bn_wordpress
  metrics:
    annotations:
      prometheus.io/port: "9104"
      prometheus.io/scrape: "true"
    enabled: false
    image: prom/mysqld-exporter
    imagePullPolicy: IfNotPresent
    imageTag: v0.10.0
    resources: {}
  persistence:
    accessMode: ReadWriteOnce
    enabled: true
    size: 8Gi
  resources:
    requests:
      cpu: 250m
      memory: 256Mi
  serviceType: ClusterIP
nodeSelector: {}
persistence:
  accessMode: ReadWriteOnce
  enabled: true
  size: 10Gi
resources:
  requests:
    cpu: 300m
    memory: 512Mi
serviceType: NodePort
wordpressBlogName: User's Blog!
wordpressEmail: user@example.com
wordpressFirstName: FirstName
wordpressLastName: LastName
wordpressUsername: user

HOOKS:
---
# wordpress-credentials-test
apiVersion: v1
kind: Pod
metadata:
  name: "wordpress-credentials-test"
  annotations:
    "helm.sh/hook": test-success
spec:
  containers:
  - name: wordpress-credentials-test
    image: bitnami/wordpress:4.9.0-r0
    env:
      - name: MARIADB_HOST
        value: wordpress-mariadb
      - name: MARIADB_PORT
        value: "3306"
      - name: WORDPRESS_DATABASE_NAME
        value: "bitnami_wordpress"
      - name: WORDPRESS_DATABASE_USER
        value: "bn_wordpress"
      - name: WORDPRESS_DATABASE_PASSWORD
        valueFrom:
          secretKeyRef:
            name: wordpress-mariadb
            key: mariadb-password
    command: ["sh", "-c", "mysql --host=$MARIADB_HOST --port=$MARIADB_PORT --user=$WORDPRESS_DATABASE_USER --password=$WORDPRESS_DATABASE_PASSWORD"]
  restartPolicy: Never
MANIFEST:

---
# Source: wordpress/charts/mariadb/templates/secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: wordpress-mariadb
  labels:
    app: wordpress-mariadb
    chart: "mariadb-0.7.0"
    release: "wordpress"
    heritage: "Tiller"
type: Opaque
data:
  mariadb-root-password: ""
  mariadb-password: ""
---
# Source: wordpress/templates/secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: wordpress-wordpress
  labels:
    app: wordpress-wordpress
    chart: "wordpress-0.7.2"
    release: "wordpress"
    heritage: "Tiller"
type: Opaque
data:

  wordpress-password: "QmQxdWpudUZvWQ=="

  smtp-password: ""
---
# Source: wordpress/charts/mariadb/templates/configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: wordpress-mariadb
  labels:
    app: wordpress-mariadb
    chart: "mariadb-0.7.0"
    release: "wordpress"
    heritage: "Tiller"
data:
  my.cnf: |-
---
# Source: wordpress/charts/mariadb/templates/pvc.yaml
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: wordpress-mariadb
  labels:
    app: wordpress-mariadb
    chart: "mariadb-0.7.0"
    release: "wordpress"
    heritage: "Tiller"
  annotations:
    volume.alpha.kubernetes.io/storage-class: default
spec:
  accessModes:
    - "ReadWriteOnce"
  resources:
    requests:
      storage: "8Gi"
---
# Source: wordpress/templates/pvc.yaml
kind: PersistentVolumeClaim
apiVersion: v1
metadata:
  name: wordpress-wordpress
  labels:
    app: wordpress-wordpress
    chart: "wordpress-0.7.2"
    release: "wordpress"
    heritage: "Tiller"
spec:
  accessModes:
    - "ReadWriteOnce"
  resources:
    requests:
      storage: "10Gi"
---
# Source: wordpress/charts/mariadb/templates/svc.yaml
apiVersion: v1
kind: Service
metadata:
  name: wordpress-mariadb
  labels:
    app: wordpress-mariadb
    chart: "mariadb-0.7.0"
    release: "wordpress"
    heritage: "Tiller"
spec:
  type: ClusterIP
  ports:
  - name: mysql
    port: 3306
    targetPort: mysql
  selector:
    app: wordpress-mariadb
---
# Source: wordpress/templates/svc.yaml
apiVersion: v1
kind: Service
metadata:
  name: wordpress-wordpress
  labels:
    app: wordpress-wordpress
    chart: "wordpress-0.7.2"
    release: "wordpress"
    heritage: "Tiller"
spec:
  type: NodePort
  ports:
  - name: http
    port: 80
    targetPort: http
  - name: https
    port: 443
    targetPort: https
  selector:
    app: wordpress-wordpress
---
# Source: wordpress/charts/mariadb/templates/deployment.yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: wordpress-mariadb
  labels:
    app: wordpress-mariadb
    chart: "mariadb-0.7.0"
    release: "wordpress"
    heritage: "Tiller"
spec:
  template:
    metadata:
      labels:
        app: wordpress-mariadb
      annotations:
        pod.alpha.kubernetes.io/init-containers: '
          [
            {
              "name": "copy-custom-config",
              "image": "bitnami/mariadb:10.1.23-r2",
              "imagePullPolicy": "IfNotPresent",
              "command": ["sh", "-c", "mkdir -p /bitnami/mariadb/conf && cp -n /bitnami/mariadb_config/my.cnf /bitnami/mariadb/conf/my_custom.cnf"],
              "volumeMounts": [
                {
                  "name": "config",
                  "mountPath": "/bitnami/mariadb_config"
                },
                {
                  "name": "data",
                  "mountPath": "/bitnami/mariadb"
                }
              ]
            }
          ]'
    spec:
      containers:
      - name: wordpress-mariadb
        image: "bitnami/mariadb:10.1.23-r2"
        imagePullPolicy: "IfNotPresent"
        env:
        - name: MARIADB_ROOT_PASSWORD
          valueFrom:
            secretKeyRef:
              name: wordpress-mariadb
              key: mariadb-root-password
        - name: MARIADB_USER
          value: "bn_wordpress"
        - name: MARIADB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: wordpress-mariadb
              key: mariadb-password
        - name: MARIADB_DATABASE
          value: "bitnami_wordpress"
        - name: ALLOW_EMPTY_PASSWORD
          value: "yes"
        ports:
        - name: mysql
          containerPort: 3306
        livenessProbe:
          exec:
            command:
            - mysqladmin
            - ping
          initialDelaySeconds: 30
          timeoutSeconds: 5
        readinessProbe:
          exec:
            command:
            - mysqladmin
            - ping
          initialDelaySeconds: 5
          timeoutSeconds: 1
        resources:
          requests:
            cpu: 250m
            memory: 256Mi

        volumeMounts:
        - name: data
          mountPath: /bitnami/mariadb
      volumes:
      - name: config
        configMap:
          name: wordpress-mariadb
      - name: data
        persistentVolumeClaim:
          claimName: wordpress-mariadb
---
# Source: wordpress/templates/deployment.yaml
apiVersion: extensions/v1beta1
kind: Deployment
metadata:
  name: wordpress-wordpress
  labels:
    app: wordpress-wordpress
    chart: "wordpress-0.7.2"
    release: "wordpress"
    heritage: "Tiller"
spec:
  replicas: 1
  template:
    metadata:
      labels:
        app: wordpress-wordpress
    spec:
      containers:
      - name: wordpress-wordpress
        image: "bitnami/wordpress:4.9.0-r0"
        imagePullPolicy: "IfNotPresent"
        env:
        - name: ALLOW_EMPTY_PASSWORD
          value: "yes"
        - name: MARIADB_ROOT_PASSWORD
          valueFrom:
            secretKeyRef:
              name: wordpress-mariadb
              key: mariadb-root-password
        - name: MARIADB_HOST
          value: wordpress-mariadb
        - name: MARIADB_PORT_NUMBER
          value: "3306"
        - name: WORDPRESS_DATABASE_NAME
          value: "bitnami_wordpress"
        - name: WORDPRESS_DATABASE_USER
          value: "bn_wordpress"
        - name: WORDPRESS_DATABASE_PASSWORD
          valueFrom:
            secretKeyRef:
              name: wordpress-mariadb
              key: mariadb-password
        - name: WORDPRESS_USERNAME
          value: "user"
        - name: WORDPRESS_PASSWORD
          valueFrom:
            secretKeyRef:
              name: wordpress-wordpress
              key: wordpress-password
        - name: WORDPRESS_EMAIL
          value: "user@example.com"
        - name: WORDPRESS_FIRST_NAME
          value: "FirstName"
        - name: WORDPRESS_LAST_NAME
          value: "LastName"
        - name: WORDPRESS_BLOG_NAME
          value: "User's Blog!"
        - name: SMTP_HOST
          value: ""
        - name: SMTP_PORT
          value: ""
        - name: SMTP_USER
          value: ""
        - name: SMTP_PASSWORD
          valueFrom:
            secretKeyRef:
              name: wordpress-wordpress
              key: smtp-password
        - name: SMTP_USERNAME
          value: ""
        - name: SMTP_PROTOCOL
          value: ""
        ports:
        - name: http
          containerPort: 80
        - name: https
          containerPort: 443
        livenessProbe:
          httpGet:
            path: /wp-login.php
            port: http
          initialDelaySeconds: 120
          timeoutSeconds: 5
          failureThreshold: 6
        readinessProbe:
          httpGet:
            path: /wp-login.php
            port: http
          initialDelaySeconds: 30
          timeoutSeconds: 3
          periodSeconds: 5
        volumeMounts:
        - mountPath: /bitnami/apache
          name: wordpress-data
          subPath: apache
        - mountPath: /bitnami/wordpress
          name: wordpress-data
          subPath: wordpress
        - mountPath: /bitnami/php
          name: wordpress-data
          subPath: php
        resources:
          requests:
            cpu: 300m
            memory: 512Mi

      volumes:
      - name: wordpress-data
        persistentVolumeClaim:
          claimName: wordpress-wordpress
```

### Step 11: Apply the HELM_CONFIG's apps (without verbose rendering).

```
$ multihelm apply
INFO[0000] Initializing MultiHelm.                       appsPath=./apps cfgFile=/tmp/hello-multihelm/configs/minikube.yaml currentContext=minikube targetContext=minikube versionNumber=v0.1.0
INFO[0000] `apply` called.
INFO[0000] Rendering app.                                app=wordpress
INFO[0000] Building chart dependencies for app.          app=wordpress chart=./src/github.com/kubernetes/charts/stable/wordpress err="<nil>"
Release "wordpress" does not exist. Installing it now.
NAME:   wordpress
LAST DEPLOYED: Thu Nov 16 20:48:25 2017
NAMESPACE: default
STATUS: DEPLOYED

RESOURCES:
==> v1/PersistentVolumeClaim
NAME                 STATUS  VOLUME                                    CAPACITY  ACCESS MODES  STORAGECLASS  AGE
wordpress-mariadb    Bound   pvc-67274581-cb39-11e7-93c7-08002744e929  8Gi       RWO           standard      1s
wordpress-wordpress  Bound   pvc-67284efe-cb39-11e7-93c7-08002744e929  10Gi      RWO           standard      1s

==> v1/Service
NAME                 TYPE       CLUSTER-IP  EXTERNAL-IP  PORT(S)                     AGE
wordpress-mariadb    ClusterIP  10.0.0.114  <none>       3306/TCP                    1s
wordpress-wordpress  NodePort   10.0.0.140  <none>       80:32543/TCP,443:30723/TCP  1s

==> v1beta1/Deployment
NAME                 DESIRED  CURRENT  UP-TO-DATE  AVAILABLE  AGE
wordpress-mariadb    1        1        1           0          1s
wordpress-wordpress  1        1        1           0          1s

==> v1/Pod(related)
NAME                                  READY  STATUS             RESTARTS  AGE
wordpress-mariadb-7cb96d9889-dmjwn    0/1    ContainerCreating  0         1s
wordpress-wordpress-7f9949cbb4-c8xmh  0/1    ContainerCreating  0         1s

==> v1/Secret
NAME                 TYPE    DATA  AGE
wordpress-mariadb    Opaque  2     1s
wordpress-wordpress  Opaque  2     1s

==> v1/ConfigMap
NAME               DATA  AGE
wordpress-mariadb  1     1s


NOTES:
1. Get the WordPress URL:

  Or running:

  export NODE_PORT=$(kubectl get --namespace default -o jsonpath="{.spec.ports[0].nodePort}" services wordpress-wordpress)
  export NODE_IP=$(kubectl get nodes --namespace default -o jsonpath="{.items[0].status.addresses[0].address}")
  echo http://$NODE_IP:$NODE_PORT/admin

2. Login with the following credentials to see your blog

  echo Username: user
  echo Password: $(kubectl get secret --namespace default wordpress-wordpress -o jsonpath="{.data.wordpress-password}" | base64 --decode)
```

### Step 12: Commit your changes to git.

```
git add .
git commit -a -m "Second commit."
```

### Step 13: Learn more at the MultiHelm readme.

https://***REMOVED***

Thanks for reading!

Feel free to email me with any comments, questions, PRs or requests.

-Joshua M. Dotson <josdotso@cisco.com>
