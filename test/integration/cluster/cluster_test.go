// +build integration

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

package cluster

import (
	"context"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func init() {
	if err := clusterv1.AddToScheme(scheme.Scheme); err != nil {
		panic(err)
	}
}

var clusterSpec = &clusterv1.ClusterSpec{
	ClusterNetwork: &clusterv1.ClusterNetwork{
		ServiceDomain: "mydomain.com",
		Pods: &clusterv1.NetworkRanges{
			CIDRBlocks: []string{"192.168.0.0/16"},
		},
	},
}

// Timeout for waiting events in seconds
const TIMEOUT = 60

func TestCluster(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Cluster-Controller")
}

var _ = Describe("Cluster-Controller", func() {
	var k8sClient *kubernetes.Clientset
	var apiclient client.Client
	var testNamespace string

	BeforeEach(func() {
		// Load configuration
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		configOverrides := &clientcmd.ConfigOverrides{}
		kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
		config, err := kubeConfig.ClientConfig()
		Expect(err).ShouldNot(HaveOccurred())

		// Create kubernetes client
		k8sClient, err = kubernetes.NewForConfig(config)
		Expect(err).ShouldNot(HaveOccurred())

		// Create namespace for test
		ns := &corev1.Namespace{ObjectMeta: metav1.ObjectMeta{GenerateName: "clusterapi-test-"}}
		ns, err = k8sClient.CoreV1().Namespaces().Create(ns)
		Expect(err).ShouldNot(HaveOccurred())
		testNamespace = ns.ObjectMeta.Name

		// Create a new client
		apiclient, err = client.New(config, client.Options{Scheme: scheme.Scheme})
		Expect(err).ShouldNot(HaveOccurred())
	})

	AfterEach(func() {
		Expect(k8sClient.CoreV1().Namespaces().Delete(testNamespace, &metav1.DeleteOptions{})).To(Succeed())
	})

	Describe("Create Cluster", func() {
		It("Should reach to pending phase after creation", func(done Done) {
			ctx := context.Background()
			// Create Cluster
			cluster := &clusterv1.Cluster{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: "cluster-",
					Namespace:    testNamespace,
				},
				Spec: *clusterSpec.DeepCopy(),
			}

			Expect(apiclient.Create(ctx, cluster)).To(Succeed())
			Eventually(func() bool {
				if err := apiclient.Get(ctx, client.ObjectKey{Name: cluster.Name, Namespace: cluster.Namespace}, cluster); err != nil {
					return false
				}
				return cluster.Status.Phase == string(clusterv1.ClusterPhasePending)
			}).Should(BeTrue())

			close(done)
		}, TIMEOUT)
	})
})
