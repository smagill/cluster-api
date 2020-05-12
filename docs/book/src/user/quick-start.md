# Quick Start

In this tutorial we'll cover the basics of how to use Cluster API to create one or more Kubernetes clusters.

## Installation

### Common Prerequisites

- Install and setup [kubectl] in your local environment
- Install [Kind] and [Docker]

### Install and/or configure a kubernetes cluster

Cluster API requires an existing Kubernetes cluster accessible via kubectl; during the installation process the
Kubernetes cluster will be transformed into a [management cluster] by installing the Cluster API [provider components], so it
is recommended to keep it separated from any application workload.

It is a common practice to create a temporary, local bootstrap cluster which is then used to provision
a target [management cluster] on the selected [infrastructure provider].

Choose one of the options below:

1. **Existing Management Cluster**

For production use-cases a "real" Kubernetes cluster should be used with appropriate backup and DR policies and procedures in place. The Kubernetes cluster must be at least v1.16+.

```bash
export KUBECONFIG=<...>
```

2. **Kind**

<aside class="note warning">

<h1>Warning</h1>

[kind] is not designed for production use.

**Minimum [kind] supported version**: v0.7.0

</aside>

[kind] can be used for creating a local Kubernetes cluster for development environments or for
the creation of a temporary [bootstrap cluster] used to provision a target [management cluster] on the selected infrastructure provider.

  ```bash
  kind create cluster
  ```
Test to ensure the local kind cluster is ready:
```
kubectl cluster-info
```

### Install clusterctl
The clusterctl CLI tool handles the lifecycle of a Cluster API management cluster.

{{#tabs name:"install-clusterctl" tabs:"linux,macOS"}}
{{#tab linux}}

#### Install clusterctl binary with curl on linux
Download the latest release; for example, to download version v0.3.0 on linux, type:
```
curl -L {{#releaselink gomodule:"sigs.k8s.io/cluster-api" asset:"clusterctl-linux-amd64" version:"0.3.x"}} -o clusterctl
```
Make the kubectl binary executable.
```
chmod +x ./clusterctl
```
Move the binary in to your PATH.
```
sudo mv ./clusterctl /usr/local/bin/clusterctl
```
Test to ensure the version you installed is up-to-date:
```
clusterctl version
```

{{#/tab }}
{{#tab macOS}}

##### Install clusterctl binary with curl on macOS
Download the latest release; for example, to download version v0.3.0 on macOS, type:
```
curl -L {{#releaselink gomodule:"sigs.k8s.io/cluster-api" asset:"clusterctl-darwin-amd64" version:"0.3.x"}} -o clusterctl
```
Make the kubectl binary executable.
```
chmod +x ./clusterctl
```
Move the binary in to your PATH.
```
sudo mv ./clusterctl /usr/local/bin/clusterctl
```
Test to ensure the version you installed is up-to-date:
```
clusterctl version
```
{{#/tab }}
{{#/tabs }}


### Initialize the management cluster

Now that we've got clusterctl installed and all the prerequisites in place, let's transform the Kubernetes cluster
into a management cluster by using `clusterctl init`.

The command accepts as input a list of providers to install; when executed for the first time, `clusterctl init`
automatically adds to the list the `cluster-api` core provider, and if unspecified, it also adds the `kubeadm` bootstrap
and `kubeadm` control-plane providers.

<aside class="note warning">

<h1>Action Required</h1>

If the  provider expects some environment variables, you should ensure those variables are set in advance. See below for
the expected settings for common providers.

Throughout this quickstart guide, we've given instructions on setting parameters using environment variables. For most
environment variables in the rest of the guide, you can also set them in ~/.cluster-api/clusterctl.yaml

See [`clusterctl init`](../clusterctl/commands/init.md) for more details.

</aside>

#### Initialization for common providers

Depending on the infrastructure provider you are planning to use, some additional prerequisites should be satisfied
before getting started with Cluster API.

{{#tabs name:"tab-installation-infrastructure" tabs:"AWS,Azure,Docker,GCP,vSphere,OpenStack,Metal3"}}
{{#tab AWS}}

Download the latest binary of `clusterawsadm` from the [AWS provider releases] and make sure to place it in your path.

The clusterawsadm command line utility assists with identity and access management (IAM) for Cluster API Provider AWS.

```bash
export AWS_REGION=us-east-1 # This is used to help encode your environment variables
export AWS_ACCESS_KEY_ID=<your-access-key>
export AWS_SECRET_ACCESS_KEY=<your-secret-access-key>
export AWS_SESSION_TOKEN=<session-token> # If you are using Multi-Factor Auth.

# The clusterawsadm utility takes the credentials that you set as environment
# variables and uses them to create a CloudFormation stack in your AWS account
# with the correct IAM resources.
clusterawsadm alpha bootstrap create-stack

# Create the base64 encoded credentials using clusterawsadm.
# This command uses your environment variables and encodes
# them in a value to be stored in a Kubernetes Secret.
export AWS_B64ENCODED_CREDENTIALS=$(clusterawsadm alpha bootstrap encode-aws-credentials)

# Finally, initialize the management cluster
clusterctl init --infrastructure aws
```

See the [AWS provider prerequisites] document for more details.

{{#/tab }}
{{#tab Azure}}

For more information about authorization, AAD, or requirements for Azure, visit the [Azure provider prerequisites] document.

```bash
export AZURE_SUBSCRIPTION_ID=<SubscriptionId>

# Create an Azure Service Principal and paste the output here
export AZURE_TENANT_ID=<Tenant>
export AZURE_CLIENT_ID=<AppId>
export AZURE_CLIENT_SECRET=<Password>

export AZURE_SUBSCRIPTION_ID_B64="$(echo -n "$AZURE_SUBSCRIPTION_ID" | base64 | tr -d '\n')"
export AZURE_TENANT_ID_B64="$(echo -n "$AZURE_TENANT_ID" | base64 | tr -d '\n')"
export AZURE_CLIENT_ID_B64="$(echo -n "$AZURE_CLIENT_ID" | base64 | tr -d '\n')"
export AZURE_CLIENT_SECRET_B64="$(echo -n "$AZURE_CLIENT_SECRET" | base64 | tr -d '\n')"

# Finally, initialize the management cluster
clusterctl init --infrastructure azure
```

{{#/tab }}
{{#tab Docker}}

If you are planning to use to test locally Cluster API using the Docker infrastructure provider, please follow additional
steps described in the [developer instruction][docker-provider] page.

{{#/tab }}
{{#tab GCP}}

```bash
# Create the base64 encoded credentials by catting your credentials json.
# This command uses your environment variables and encodes
# them in a value to be stored in a Kubernetes Secret.
export GCP_B64ENCODED_CREDENTIALS=$( cat /path/to/gcp-credentials.json | base64 | tr -d '\n' )

# Finally, initialize the management cluster
clusterctl init --infrastructure gcp
```

{{#/tab }}
{{#tab vSphere}}

```bash
# The username used to access the remote vSphere endpoint
export VSPHERE_USERNAME="vi-admin@vsphere.local"
# The password used to access the remote vSphere endpoint
# You may want to set this in ~/.cluster-api/clusterctl.yaml so your password is not in
# bash history
export VSPHERE_PASSWORD="admin!23"

# Finally, initialize the management cluster
clusterctl init --infrastructure vsphere
```

For more information about prerequisites, credentials management, or permissions for vSphere, see the [vSphere
project][vSphere getting started guide].

{{#/tab }}
{{#tab OpenStack}}

```bash
# Initialize the management cluster
clusterctl init --infrastructure openstack
```

{{#/tab }}
{{#tab Metal3}}

Please visit the [Metal3 project][Metal3 getting started guide].

{{#/tab }}
{{#/tabs }}


The output of `clusterctl init` is similar to this:

```shell
Fetching providers
Installing cert-manager
Waiting for cert-manager to be available...
Installing Provider="cluster-api" Version="v0.3.0" TargetNamespace="capi-system"
Installing Provider="bootstrap-kubeadm" Version="v0.3.0" TargetNamespace="capi-kubeadm-bootstrap-system"
Installing Provider="control-plane-kubeadm" Version="v0.3.0" TargetNamespace="capi-kubeadm-control-plane-system"
Installing Provider="infrastructure-aws" Version="v0.5.0" TargetNamespace="capa-system"

Your management cluster has been initialized successfully!

You can now create your first workload cluster by running the following:

  clusterctl config cluster [name] --kubernetes-version [version] | kubectl apply -f -
```

### Create your first workload cluster

Once the management cluster is ready, you can create your first workload cluster.

#### Preparing the workload cluster configuration

The `clusterctl config cluster` command returns a YAML template for creating a [workload cluster].

<aside class="note">

<h1> Which provider will be used for my cluster? </h1>

The `clusterctl config cluster` command uses smart defaults in order to simplify the user experience; in this example,
it detects that there is only an `aws` infrastructure provider and so it uses that when creating the cluster.

</aside>

<aside class="note">

<h1> What topology will be used for my cluster? </h1>

The `clusterctl config cluster` command by default uses cluster templates which are provided by the infrastructure
providers. See the provider's documentation for more information.

See the `clusterctl config cluster` [command][clusterctl config cluster] documentation for
details about how to use alternative sources. for cluster templates.

</aside>

<aside class="note warning">

<h1>Action Required</h1>

If the cluster template defined by the infrastructure provider expects some environment variables, you
should ensure those variables are set in advance.

Instructions are provided for common providers below.

Otherwise, you can look at the `clusterctl config cluster` [command][clusterctl config cluster] documentation for details about how to
discover the list of variables required by a cluster templates.

</aside>

#### Required configuration for common providers

Depending on the infrastructure provider you are planning to use, some additional prerequisites should be satisfied
before configuring a cluster with Cluster API.

{{#tabs name:"tab-configuration-infrastructure" tabs:"AWS,Azure,Docker,GCP,vSphere,OpenStack,Metal3"}}
{{#tab AWS}}

Download the latest binary of `clusterawsadm` from the [AWS provider releases] and make sure to place it in your path.

```bash
export AWS_REGION=us-east-1
export AWS_SSH_KEY_NAME=default
# Select instance types
export AWS_CONTROL_PLANE_MACHINE_TYPE=t3.large
export AWS_NODE_MACHINE_TYPE=t3.large
```

See the [AWS provider prerequisites] document for more details.

{{#/tab }}
{{#tab Azure}}

```bash
export CLUSTER_NAME="capi-quickstart"
# Name of the virtual network in which to provision the cluster.
export AZURE_VNET_NAME=${CLUSTER_NAME}-vnet
# Name of the resource group to provision into
export AZURE_RESOURCE_GROUP=${CLUSTER_NAME}
# Name of the Azure datacenter location
export AZURE_LOCATION="centralus"
# Select machine types
export AZURE_CONTROL_PLANE_MACHINE_TYPE="Standard_D2s_v3"
export AZURE_NODE_MACHINE_TYPE="Standard_D2s_v3"
# Generate SSH key.
# If you want to provide your own key, skip this step and set AZURE_SSH_PUBLIC_KEY to your existing public key.
SSH_KEY_FILE=.sshkey
rm -f "${SSH_KEY_FILE}" 2>/dev/null
ssh-keygen -t rsa -b 2048 -f "${SSH_KEY_FILE}" -N '' 1>/dev/null
echo "Machine SSH key generated in ${SSH_KEY_FILE}"
export AZURE_SSH_PUBLIC_KEY=$(cat "${SSH_KEY_FILE}.pub" | base64 | tr -d '\r\n')
```

For more information about authorization, AAD, or requirements for Azure, visit the [Azure provider prerequisites] document.

{{#/tab }}
{{#tab Docker}}

If you are planning to use to test locally Cluster API using the Docker infrastructure provider, please follow additional
steps described in the [developer instructions][docker-provider] page.

{{#/tab }}
{{#tab GCP}}

See the [GCP provider] for more information.

{{#/tab }}
{{#tab vSphere}}

It is required to use an official CAPV machine images for your vSphere VM templates. See [uploading CAPV machine images][capv-upload-images] for instructions on how to do this.

```bash
# The vCenter server IP or FQDN
export VSPHERE_SERVER="10.0.0.1"
# The vSphere datacenter to deploy the management cluster on
export VSPHERE_DATACENTER="SDDC-Datacenter"
# The vSphere datastore to deploy the management cluster on
export VSPHERE_DATASTORE="vsanDatastore"
# The VM network to deploy the management cluster on
export VSPHERE_NETWORK="VM Network"
# The vSphere resource pool for your VMs
export VSPHERE_RESOURCE_POOL="*/Resources"
# The VM folder for your VMs. Set to "" to use the root vSphere folder
export VSPHERE_FOLDER="vm"
# The VM template to use for your VMs
export VSPHERE_TEMPLATE="ubuntu-1804-kube-v1.17.3"
# The VM template to use for the HAProxy load balancer of the management cluster
export VSPHERE_HAPROXY_TEMPLATE="capv-haproxy-v0.6.0-rc.2"
# The public ssh authorized key on all machines
export VSPHERE_SSH_AUTHORIZED_KEY="ssh-rsa AAAAB3N..."

clusterctl init --infrastructure vsphere
```

For more information about prerequisites, credentials management, or permissions for vSphere, see the [vSphere getting started guide].

{{#/tab }}
{{#tab OpenStack}}

A ClusterAPI compatible image must be available in your OpenStack. For instructions on how to build a compatible image
see [image-builder](https://image-builder.sigs.k8s.io/capi/capi.html).
Depending on your OpenStack and underlying hypervisor the following options might be of interest:
* [image-builder (OpenStack)](https://image-builder.sigs.k8s.io/capi/providers/openstack.html)
* [image-builder (vSphere)](https://image-builder.sigs.k8s.io/capi/providers/vsphere.html)

To see all required OpenStack environment variables execute:
```bash
clusterctl config cluster --infrastructure openstack --list-variables capi-quickstart
```

The following script can be used to export some of them:
```bash
wget https://raw.githubusercontent.com/kubernetes-sigs/cluster-api-provider-openstack/master/templates/env.rc -O /tmp/env.rc
source /tmp/env.rc <path/to/clouds.yaml> <cloud>
```

A full configuration reference can be found in [configuration.md](https://github.com/kubernetes-sigs/cluster-api-provider-openstack/blob/master/docs/configuration.md).

{{#/tab }}
{{#tab Metal3}}

Please visit the [Metal3 getting started guide].

{{#/tab }}
{{#/tabs }}

#### Generating the cluster configuration

For the purpose of this tutorial, we'll name our cluster capi-quickstart.

```bash
clusterctl config cluster capi-quickstart --kubernetes-version v1.17.3 --control-plane-machine-count=3 --worker-machine-count=3 > capi-quickstart.yaml
```

This creates a YAML file named `capi-quickstart.yaml` with a predefined list of Cluster API objects; Cluster, Machines,
Machine Deployments, etc.

The file can be eventually modified using your editor of choice.

See `[clusterctl config cluster]` for more details.

#### Apply the workload cluster

When ready, run the following command to apply the cluster manifest.

```bash
kubectl apply -f capi-quickstart.yaml
```

The output is similar to this:

```bash
cluster.cluster.x-k8s.io/capi-quickstart created
awscluster.infrastructure.cluster.x-k8s.io/capi-quickstart created
kubeadmcontrolplane.controlplane.cluster.x-k8s.io/capi-quickstart-control-plane created
awsmachinetemplate.infrastructure.cluster.x-k8s.io/capi-quickstart-control-plane created
machinedeployment.cluster.x-k8s.io/capi-quickstart-md-0 created
awsmachinetemplate.infrastructure.cluster.x-k8s.io/capi-quickstart-md-0 created
kubeadmconfigtemplate.bootstrap.cluster.x-k8s.io/capi-quickstart-md-0 created
```

#### Accessing the workload cluster

The cluster will now start provisioning. You can check status with:

```bash
kubectl get cluster --all-namespaces
```

To verify the first control plane is up:

```bash
kubectl get kubeadmcontrolplane --all-namespaces
```

You should see an output is similar to this:

```bash
NAME                              READY   INITIALIZED   REPLICAS   READY REPLICAS   UPDATED REPLICAS   UNAVAILABLE REPLICAS
capi-quickstart-control-plane             true          3                           3                  3
```

<aside class="note warning">

<h1> Warning </h1>

The control planes won't be `Ready` until we install a CNI in the next step.

</aside>

After the first control plane node is up and running, we can retrieve the [workload cluster] Kubeconfig:

```bash
kubectl --namespace=default get secret/capi-quickstart-kubeconfig -o jsonpath={.data.value} \
  | base64 --decode \
  > ./capi-quickstart.kubeconfig
```

### Deploy a CNI solution

Calico is used here as an example.

{{#tabs name:"tab-deploy-cni" tabs:"AWS|Docker|GCP|vSphere|OpenStack|Metal3,Azure"}}
{{#tab AWS|Docker|GCP|vSphere|OpenStack|Metal3}}

```bash
kubectl --kubeconfig=./capi-quickstart.kubeconfig \
  apply -f https://docs.projectcalico.org/v3.12/manifests/calico.yaml
```

After a short while, our nodes should be running and in `Ready` state,
let's check the status using `kubectl get nodes`:

```bash
kubectl --kubeconfig=./capi-quickstart.kubeconfig get nodes
```

{{#/tab }}
{{#tab Azure}}

Azure [does not currently support Calico networking](https://docs.projectcalico.org/reference/public-cloud/azure). As a workaround, it is recommended that Azure clusters use the Calico spec below that uses VXLAN.

```bash
kubectl --kubeconfig=./capi-quickstart.kubeconfig \
  apply -f https://raw.githubusercontent.com/kubernetes-sigs/cluster-api-provider-azure/master/templates/addons/calico.yaml
```

After a short while, our nodes should be running and in `Ready` state,
let's check the status using `kubectl get nodes`:

```bash
kubectl --kubeconfig=./capi-quickstart.kubeconfig get nodes
```

{{#/tab }}
{{#/tabs }}

## Next steps

See the [clusterctl] documentation for more detail about clusterctl supported actions.

<!-- links -->
[AWS provider prerequisites]: https://github.com/kubernetes-sigs/cluster-api-provider-aws/blob/master/docs/prerequisites.md
[AWS provider releases]: https://github.com/kubernetes-sigs/cluster-api-provider-aws/releases
[Azure Provider Prerequisites]: https://github.com/kubernetes-sigs/cluster-api-provider-azure/blob/master/docs/getting-started.md#prerequisites
[bootstrap cluster]: ../reference/glossary.md#bootstrap-cluster
[capv-upload-images]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/blob/master/docs/getting_started.md#uploading-the-machine-images
[clusterctl config cluster]: ../clusterctl/commands/config-cluster.md
[clusterctl]: ../clusterctl/overview.md
[Docker]: https://www.docker.com/
[docker-provider]: ../clusterctl/developers.md#additional-steps-in-order-to-use-the-docker-provider
[GCP provider]: https://github.com/kubernetes-sigs/cluster-api-provider-gcp
[infrastructure provider]: ../reference/glossary.md#infrastructure-provider
[kind]: https://kind.sigs.k8s.io/
[KubeadmControlPlane]: ../developer/architecture/controllers/control-plane.md
[kubectl]: https://kubernetes.io/docs/tasks/tools/install-kubectl/
[management cluster]: ../reference/glossary.md#management-cluster
[Metal3 getting started guide]: https://github.com/metal3-io/cluster-api-provider-metal3/
[provider components]: ../reference/glossary.md#provider-components
[vSphere getting started guide]: https://github.com/kubernetes-sigs/cluster-api-provider-vsphere/
[workload cluster]: ../reference/glossary.md#workload-cluster
