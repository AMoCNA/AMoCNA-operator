# AMoCNA Operator
AMoCNA Operator aims to simplify usage of [AMoCNA](https://www.researchgate.net/publication/344415012_Autonomic_Management_Framework_for_Cloud-Native_Applications) framework provaiding a Clound - Native applications autonomy. 

Currently AMoCNA Operator provides utilities that greatly aid AMoCNA Ddeployment with use of a [Hephaestus AMoCNA implementation](https://github.com/Hephaestus-Metrics). Once user creates a rule evaluation Image (referred to as a Metrics Adapter), execution controller and deploys Prometheus, all that has to be done is creating an appropriate Custom Resource and applying it to the envirnoment. Sample components that need to be implemented by a user can be found in Hephaestus Project Repository. 

This implementaion of AMoCNA operator was build with a help of a [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder).

## Custom Resource Definition
Currently AMoCNA Operator uses Custom Resource Definiton which follows the below definition:
```
(Required)

// Docker tag of hephaestusmetrics/gui image to use
HephaestusGuiVersion: string

// Address on which Prometheus is exposed on cluster
PrometheusAddress: string

// Docker image of an Execution Controller
ExecutionControllerImage: string

// Docker image of a Metrics Adapter
MetricsAdapterImage: string


(Optional)

// YAML containig Config Map used by Hephaestus GUI, defaults to empty collection
HephaestusGuiConfigMapRaw: map[string]string

// Container Port on which Metrics Adapter app is exposed, defaults to 8085
MetricsAdapterInternalPort: int32

// Container Port on which Hephaestus GUI app is exposed, defaults to 8080
HephaestusGuiInternalPort: int32

// Node Port on which Hephaestus GUI service is exposed, defaults to 31122
HephaestusGuiExternalPort: int32

// Container Port on which Hephaestus Execution Controller app is exposed, defaults to 8097
ExecutionControllerInternalPort: int32

// Service Account name provided to an Execution Controller
ExecutionControllerServiceAccountName: string


(Currently not supported)

// File path of a YAML containig Config Map used by Hephaestus GUI. 
// Currently onlyHephaestusGuiConfigMapRaw is supported instead
HephaestusGuiConfigMapFilePath: string
```


## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running on the cluster
1. Install Instances of Custom Resources:

```sh
kubectl apply -f config/samples/
```

2. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/operator:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/operator:tag
```


### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

#### Deployment
The operator deploys AMoCNA components in the following order:
* PVC for the volume used by the Hephaestus GUI
* Config Map used by the Hephaestus GUI
* Hephaestus GUI
* Hephaestus GUI Service
* Metrics Adapter
* Execution Controller

### Test It Out
1. Install the CRDs into the cluster:

```sh
make install
```

2. Run your controller (this will run in the foreground, so switch to a new terminal if you want to leave it running):

```sh
make run
```

**NOTE:** You can also run this in one step by running: `make install run`

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

