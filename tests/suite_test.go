/*
Copyright 2019 The OpenEBS Authors
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

package tests

import (
	"context"
	"flag"
	"os"
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ns "github.com/openebs/maya/pkg/kubernetes/namespace/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	// auth plugins
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"

	"github.com/openebs/dynamic-localpv-provisioner/tests/disk"
)

var (
	kubeConfigPath                  string
	openebsNamespace                string
	namespace                       = "localpv-ns"
	namespaceObj                    *corev1.Namespace
	storageClassLabelSelector       = "openebs.io/test-sc=true"
	err                             error
	LocalPVProvisionerLabelSelector = "openebs.io/component-name=openebs-localpv-provisioner"
	hostpathDir                     = "/var/openebs/integration-test"
	xfsHostpathDir                  = "/var/openebs/integration-test/xfs/"
	physicalDisk                    = disk.Disk{}
)

func TestSource(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test application deployment")
}

func init() {
	flag.StringVar(&kubeConfigPath, "kubeconfig", os.Getenv("KUBECONFIG"), "path to kubeconfig to invoke kubernetes API calls")
	flag.StringVar(&openebsNamespace, "openebs-namespace", "openebs", "kubernetes namespace where the OpenEBS components are present")
}

var ops *Operations

var _ = BeforeSuite(func() {

	ops = NewOperations(WithKubeConfigPath(kubeConfigPath))

	By("waiting for openebs-localpv-provisioner pod to come into running state")
	provPodCount := ops.GetPodRunningCountEventually(
		openebsNamespace,
		LocalPVProvisionerLabelSelector,
		1,
	)
	Expect(provPodCount).To(Equal(1))

	By("building a namespace")
	namespaceObj, err = ns.NewBuilder().
		WithGenerateName(namespace).
		APIObject()
	Expect(err).ShouldNot(HaveOccurred(), "while building namespace {%s}", namespaceObj.GenerateName)

	By("creating above namespace")
	namespaceObj, err = ops.NSClient.Create(namespaceObj)
	Expect(err).To(BeNil(), "while creating namespace {%s}", namespaceObj.GenerateName)

	By("preparing the loopback device with xfs fs")
	physicalDisk, err = disk.PrepareDisk("xfs", xfsHostpathDir)
	Expect(err).To(BeNil(), "while preparing disk {%+v}", physicalDisk)
})

var _ = AfterSuite(func() {

	By("deleting namespace")
	err = ops.NSClient.Delete(namespaceObj.Name, &metav1.DeleteOptions{})
	Expect(err).To(BeNil(), "while deleting namespace {%s}", namespaceObj.Name)

	By("deleting test StorageClasses")
	err = ops.SCClient.DeleteCollection(
		context.TODO(),
		metav1.ListOptions{
			LabelSelector: storageClassLabelSelector,
		},
		&metav1.DeleteOptions{},
	)

	By("destroying the created disk")
	err = physicalDisk.DestroyDisk(xfsHostpathDir)
	Expect(err).To(BeNil(), "while destroying the disk {%+v}", physicalDisk)
})
