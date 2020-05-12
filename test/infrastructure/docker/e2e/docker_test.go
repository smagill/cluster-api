// +build e2e

/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package e2e

import (
	"fmt"
	"net/http"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/runtime"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	bootstrapv1 "sigs.k8s.io/cluster-api/bootstrap/kubeadm/api/v1alpha3"
	"sigs.k8s.io/cluster-api/bootstrap/kubeadm/types/v1beta1"
	controlplanev1 "sigs.k8s.io/cluster-api/controlplane/kubeadm/api/v1alpha3"
	"sigs.k8s.io/cluster-api/test/framework"
	infrav1 "sigs.k8s.io/cluster-api/test/infrastructure/docker/api/v1alpha3"
	"sigs.k8s.io/cluster-api/util"
	ctrlclient "sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Docker Create", func() {
	var (
		namespace  = "default"
		clusterGen = newClusterGenerator("create")
		mgmtClient ctrlclient.Client
		cluster    *clusterv1.Cluster
	)
	SetDefaultEventuallyTimeout(15 * time.Minute)
	SetDefaultEventuallyPollingInterval(10 * time.Second)

	AfterEach(func() {
		// Delete the workload cluster
		deleteClusterInput := framework.DeleteClusterInput{
			Deleter: mgmtClient,
			Cluster: cluster,
		}
		framework.DeleteCluster(ctx, deleteClusterInput)

		waitForClusterDeletedInput := framework.WaitForClusterDeletedInput{
			Getter:  mgmtClient,
			Cluster: cluster,
		}
		framework.WaitForClusterDeleted(ctx, waitForClusterDeletedInput)

		assertAllClusterAPIResourcesAreGoneInput := framework.AssertAllClusterAPIResourcesAreGoneInput{
			Lister:  mgmtClient,
			Cluster: cluster,
		}
		framework.AssertAllClusterAPIResourcesAreGone(ctx, assertAllClusterAPIResourcesAreGoneInput)

		ensureDockerDeletedInput := ensureDockerArtifactsDeletedInput{
			Lister:  mgmtClient,
			Cluster: cluster,
		}
		ensureDockerArtifactsDeleted(ensureDockerDeletedInput)

		// Dump cluster API and docker related resources to artifacts before deleting them.
		Expect(framework.DumpResources(mgmt, resourcesPath, GinkgoWriter)).To(Succeed())
		resources := map[string]runtime.Object{
			"DockerCluster":         &infrav1.DockerClusterList{},
			"DockerMachine":         &infrav1.DockerMachineList{},
			"DockerMachineTemplate": &infrav1.DockerMachineTemplateList{},
		}
		Expect(framework.DumpProviderResources(mgmt, resources, resourcesPath, GinkgoWriter)).To(Succeed())
	})

	Specify("multi-node cluster with failure domains", func() {
		replicas := 3
		var (
			infraCluster *infrav1.DockerCluster
			template     *infrav1.DockerMachineTemplate
			controlPlane *controlplanev1.KubeadmControlPlane
			err          error
		)
		cluster, infraCluster, controlPlane, template = clusterGen.GenerateCluster(namespace, int32(replicas))
		// Set failure domains here
		infraCluster.Spec.FailureDomains = clusterv1.FailureDomains{
			"domain-one":   {ControlPlane: true},
			"domain-two":   {ControlPlane: true},
			"domain-three": {ControlPlane: true},
			"domain-four":  {ControlPlane: false},
		}

		md, infraTemplate, bootstrapTemplate := GenerateMachineDeployment(cluster, 1)

		// Set up the client to the management cluster
		mgmtClient, err = mgmt.GetClient()
		Expect(err).NotTo(HaveOccurred())

		// Set up the cluster object
		createClusterInput := framework.CreateClusterInput{
			Creator:      mgmtClient,
			Cluster:      cluster,
			InfraCluster: infraCluster,
		}
		framework.CreateCluster(ctx, createClusterInput)

		// Set up the KubeadmControlPlane
		createKubeadmControlPlaneInput := framework.CreateKubeadmControlPlaneInput{
			Creator:         mgmtClient,
			ControlPlane:    controlPlane,
			MachineTemplate: template,
		}
		framework.CreateKubeadmControlPlane(ctx, createKubeadmControlPlaneInput)

		// Wait for the cluster to provision.
		assertClusterProvisionsInput := framework.WaitForClusterToProvisionInput{
			Getter:  mgmtClient,
			Cluster: cluster,
		}
		framework.WaitForClusterToProvision(ctx, assertClusterProvisionsInput)

		// Wait for at least one control plane node to be ready
		waitForOneKubeadmControlPlaneMachineToExistInput := framework.WaitForOneKubeadmControlPlaneMachineToExistInput{
			Lister:       mgmtClient,
			Cluster:      cluster,
			ControlPlane: controlPlane,
		}
		framework.WaitForOneKubeadmControlPlaneMachineToExist(ctx, waitForOneKubeadmControlPlaneMachineToExistInput, "15m")

		// Insatll a networking solution on the workload cluster
		workloadClient, err := mgmt.GetWorkloadClient(ctx, cluster.Namespace, cluster.Name)
		Expect(err).ToNot(HaveOccurred())
		applyYAMLURLInput := framework.ApplyYAMLURLInput{
			Client:        workloadClient,
			HTTPGetter:    http.DefaultClient,
			NetworkingURL: "https://docs.projectcalico.org/manifests/calico.yaml",
			Scheme:        mgmt.Scheme,
		}
		framework.ApplyYAMLURL(ctx, applyYAMLURLInput)

		// Wait for the controlplane nodes to exist
		assertKubeadmControlPlaneNodesExistInput := framework.WaitForKubeadmControlPlaneMachinesToExistInput{
			Lister:       mgmtClient,
			Cluster:      cluster,
			ControlPlane: controlPlane,
		}
		framework.WaitForKubeadmControlPlaneMachinesToExist(ctx, assertKubeadmControlPlaneNodesExistInput, "15m", "10s")

		// Create the workload nodes
		createMachineDeploymentinput := framework.CreateMachineDeploymentInput{
			Creator:                 mgmtClient,
			MachineDeployment:       md,
			BootstrapConfigTemplate: bootstrapTemplate,
			InfraMachineTemplate:    infraTemplate,
		}
		framework.CreateMachineDeployment(ctx, createMachineDeploymentinput)

		// Wait for the workload nodes to exist
		waitForMachineDeploymentNodesToExistInput := framework.WaitForMachineDeploymentNodesToExistInput{
			Lister:            mgmtClient,
			Cluster:           cluster,
			MachineDeployment: md,
		}
		framework.WaitForMachineDeploymentNodesToExist(ctx, waitForMachineDeploymentNodesToExistInput)

		// Wait for the control plane to be ready
		waitForControlPlaneToBeReadyInput := framework.WaitForControlPlaneToBeReadyInput{
			Getter:       mgmtClient,
			ControlPlane: controlPlane,
		}
		framework.WaitForControlPlaneToBeReady(ctx, waitForControlPlaneToBeReadyInput)

		// Assert failure domain is working as expected
		assertControlPlaneFailureDomainInput := framework.AssertControlPlaneFailureDomainsInput{
			GetLister:  mgmtClient,
			ClusterKey: util.ObjectKey(cluster),
			ExpectedFailureDomains: map[string]int{
				"domain-one":   1,
				"domain-two":   1,
				"domain-three": 1,
				"domain-four":  0,
			},
		}
		framework.AssertControlPlaneFailureDomains(ctx, assertControlPlaneFailureDomainInput)

		Describe("Docker recover from manual workload machine deletion", func() {
			By("cleaning up etcd members and kubeadm configMap")
			inClustersNamespaceListOption := ctrlclient.InNamespace(cluster.Namespace)
			// ControlPlane labels
			matchClusterListOption := ctrlclient.MatchingLabels{
				clusterv1.MachineControlPlaneLabelName: "",
				clusterv1.ClusterLabelName:             cluster.Name,
			}

			machineList := &clusterv1.MachineList{}
			err = mgmtClient.List(ctx, machineList, inClustersNamespaceListOption, matchClusterListOption)
			Expect(err).ToNot(HaveOccurred())
			Expect(machineList.Items).To(HaveLen(int(*controlPlane.Spec.Replicas)))

			Expect(mgmtClient.Delete(ctx, &machineList.Items[0])).To(Succeed())

			Eventually(func() (int, error) {
				machineList := &clusterv1.MachineList{}
				if err := mgmtClient.List(ctx, machineList, inClustersNamespaceListOption, matchClusterListOption); err != nil {
					fmt.Println(err)
					return 0, err
				}
				return len(machineList.Items), nil
			}, "15m", "5s").Should(Equal(int(*controlPlane.Spec.Replicas) - 1))

			By("ensuring a replacement machine is created")
			Eventually(func() (int, error) {
				machineList := &clusterv1.MachineList{}
				if err := mgmtClient.List(ctx, machineList, inClustersNamespaceListOption, matchClusterListOption); err != nil {
					fmt.Println(err)
					return 0, err
				}
				return len(machineList.Items), nil
			}, "15m", "30s").Should(Equal(int(*controlPlane.Spec.Replicas)))
		})
	})

})

func GenerateMachineDeployment(cluster *clusterv1.Cluster, replicas int32) (*clusterv1.MachineDeployment, *infrav1.DockerMachineTemplate, *bootstrapv1.KubeadmConfigTemplate) {
	namespace := cluster.GetNamespace()
	generatedName := fmt.Sprintf("%s-md", cluster.GetName())
	version := "1.16.3"

	infraTemplate := &infrav1.DockerMachineTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      generatedName,
		},
		Spec: infrav1.DockerMachineTemplateSpec{},
	}

	bootstrap := &bootstrapv1.KubeadmConfigTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      generatedName,
		},
	}

	template := clusterv1.MachineTemplateSpec{
		ObjectMeta: clusterv1.ObjectMeta{
			Namespace: namespace,
			Name:      generatedName,
		},
		Spec: clusterv1.MachineSpec{
			ClusterName: cluster.GetName(),
			Bootstrap: clusterv1.Bootstrap{
				ConfigRef: &corev1.ObjectReference{
					APIVersion: bootstrapv1.GroupVersion.String(),
					Kind:       framework.TypeToKind(bootstrap),
					Namespace:  bootstrap.GetNamespace(),
					Name:       bootstrap.GetName(),
				},
			},
			InfrastructureRef: corev1.ObjectReference{
				APIVersion: infrav1.GroupVersion.String(),
				Kind:       framework.TypeToKind(infraTemplate),
				Namespace:  infraTemplate.GetNamespace(),
				Name:       infraTemplate.GetName(),
			},
			Version: &version,
		},
	}

	machineDeployment := &clusterv1.MachineDeployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      generatedName,
		},
		Spec: clusterv1.MachineDeploymentSpec{
			ClusterName:             cluster.GetName(),
			Replicas:                &replicas,
			Template:                template,
			Strategy:                nil,
			MinReadySeconds:         nil,
			RevisionHistoryLimit:    nil,
			Paused:                  false,
			ProgressDeadlineSeconds: nil,
		},
	}
	return machineDeployment, infraTemplate, bootstrap
}

type clusterGenerator struct {
	prefix  string
	counter int
}

func newClusterGenerator(name string) *clusterGenerator {
	var prefix string
	if len(name) != 0 {
		prefix = fmt.Sprintf("test-%s-", name)
	} else {
		prefix = "test-"
	}

	return &clusterGenerator{
		prefix: prefix,
	}
}

func (c *clusterGenerator) GenerateCluster(namespace string, replicas int32) (*clusterv1.Cluster, *infrav1.DockerCluster, *controlplanev1.KubeadmControlPlane, *infrav1.DockerMachineTemplate) {
	generatedName := fmt.Sprintf("%s%d", c.prefix, c.counter)
	c.counter++
	version := "v1.16.3"

	infraCluster := &infrav1.DockerCluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      generatedName,
		},
	}

	template := &infrav1.DockerMachineTemplate{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      generatedName,
		},
		Spec: infrav1.DockerMachineTemplateSpec{},
	}

	kcp := &controlplanev1.KubeadmControlPlane{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      generatedName,
		},
		Spec: controlplanev1.KubeadmControlPlaneSpec{
			Replicas: &replicas,
			Version:  version,
			InfrastructureTemplate: corev1.ObjectReference{
				Kind:       framework.TypeToKind(template),
				Namespace:  template.GetNamespace(),
				Name:       template.GetName(),
				APIVersion: infrav1.GroupVersion.String(),
			},
			KubeadmConfigSpec: bootstrapv1.KubeadmConfigSpec{
				ClusterConfiguration: &v1beta1.ClusterConfiguration{
					APIServer: v1beta1.APIServer{
						// Darwin support
						CertSANs: []string{"127.0.0.1"},
					},
				},
				InitConfiguration: &v1beta1.InitConfiguration{},
				JoinConfiguration: &v1beta1.JoinConfiguration{},
			},
		},
	}

	cluster := &clusterv1.Cluster{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      generatedName,
		},
		Spec: clusterv1.ClusterSpec{
			ClusterNetwork: &clusterv1.ClusterNetwork{
				Services: &clusterv1.NetworkRanges{CIDRBlocks: []string{}},
				Pods:     &clusterv1.NetworkRanges{CIDRBlocks: []string{"192.168.0.0/16"}},
			},
			InfrastructureRef: &corev1.ObjectReference{
				APIVersion: infrav1.GroupVersion.String(),
				Kind:       framework.TypeToKind(infraCluster),
				Namespace:  infraCluster.GetNamespace(),
				Name:       infraCluster.GetName(),
			},
			ControlPlaneRef: &corev1.ObjectReference{
				APIVersion: controlplanev1.GroupVersion.String(),
				Kind:       framework.TypeToKind(kcp),
				Namespace:  kcp.GetNamespace(),
				Name:       kcp.GetName(),
			},
		},
	}
	return cluster, infraCluster, kcp, template
}
