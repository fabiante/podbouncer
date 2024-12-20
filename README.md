# podbouncer

Podbouncer is a [Kubernetes operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
responsible for deleting non-running pods after a configurable grace period.

## Description

Using podbouncer in your cluster will result in the deletion of all pods in one of
these states: `Pending`, `Completed`, `Failed`

This operator acts on pods of all namespaces, except the `kube-system` namespace.

This controller is fairly simple and has only one configuration option can configure
via a ConfigMap: The `maxPodAge` field controls the maximum age a non-running pod
may have before it will be deleted by this controller.

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: config
  labels:
    control-plane: controller-manager
    app.kubernetes.io/name: podbouncer
    app.kubernetes.io/managed-by: kustomize
data:
  maxPodAge: "1h"
```

## Quick Start

If you simply want to run podbouncer on your cluster, you can use the command below:

**Warning:** With this command you will run a potentially destructive controller in your cluster.

_Note:_ This installs podbouncer in the `podbouncer-system` namespace.

```shell
kubectl apply -f https://raw.githubusercontent.com/fabiante/podbouncer/refs/heads/main/dist/install.yaml
```

## Constraints

This operator currently requires the usage of the `podbouncer-system` namespace and will
only load its configuration from a ConfigMap object named `podbouncer-config` (in the same namespace).

## Getting Started (Local Build + Deploy)

### Prerequisites
- go version v1.22.0+
- docker version 17.03+.
- kubectl version v1.11.3+.
- Access to a Kubernetes v1.11.3+ cluster.

### To Deploy on the cluster

**Build and push your image to the location specified by `IMG`:**

```sh
make docker-build docker-push IMG=<some-registry>/podbouncer:latest
```

**NOTE:** This image ought to be published in the personal registry you specified.
And it is required to have access to pull the image from the working environment.
Make sure you have the proper permission to the registry if the above commands don’t work.

**Install the CRDs into the cluster:**

Since the controller does not use any custom resources, no CRDs have to be installed in your cluster.

**Deploy the Manager to the cluster with the image specified by `IMG`:**

```sh
make deploy IMG=<some-registry>/podbouncer:latest
```

> **NOTE**: If you encounter RBAC errors, you may need to grant yourself cluster-admin
privileges or be logged in as admin.

**Change controller configuration**
If you want to configure the controller, modify the ConfigMap in `config/samples/config.yaml`
and re-deploy the project

```sh
make deploy IMG=<some-registry>/podbouncer:latest
```

### To Uninstall

**UnDeploy the controller from the cluster:**

```sh
make undeploy
```

## Project Distribution

Following are the steps to build the installer and distribute this project to users.

1. Build the installer for the image built and published in the registry:

```sh
make build-installer IMG=<some-registry>/podbouncer:latest
```

NOTE: The makefile target mentioned above generates an 'install.yaml'
file in the dist directory. This file contains all the resources built
with Kustomize, which are necessary to install this project without
its dependencies.

2. Using the installer

Users can just run kubectl apply -f <URL for YAML BUNDLE> to install the project, i.e.:

```sh
kubectl apply -f https://raw.githubusercontent.com/<org>/podbouncer/<tag or branch>/dist/install.yaml
```

## Testing

### Create Test Pods

To test podbouncer, you must have some non-running pods in your cluster:

```shell
kubectl run --restart=Never --image busybox some-pod
kubectl run --restart=Never --image busybox some-other-pod
```

## Contributing

This project is a personal learning project.

You are welcome to use it, learn from it. You are welcome to submit pull requests,
but I do not guarantee any level of activity on this project.

## License

Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

