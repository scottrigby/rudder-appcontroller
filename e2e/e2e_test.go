// Copyright 2017 Mirantis
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

package e2e

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
)

var _ = Describe("Basic Suite", func() {

	var helm HelmManager
	var namespace *v1.Namespace
	var clientset kubernetes.Interface

	BeforeEach(func() {
		var err error
		clientset, err = KubeClient()
		Expect(err).NotTo(HaveOccurred())
		By("Creating namespace and initializing test framework")
		namespaceObj := &v1.Namespace{
			ObjectMeta: v1.ObjectMeta{
				GenerateName: "e2e-appcontroller-rudder-",
			},
		}
		namespace, err = clientset.Core().Namespaces().Create(namespaceObj)
		Expect(err).NotTo(HaveOccurred())
		helm = &BinaryHelmManager{
			Namespace: namespace.Name,
			Clientset: clientset,
			HelmBin:   "helm",
		}
		Expect(helm.InstallTiller()).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		By("Removing namespace")
		DeleteNS(clientset, namespace)
		By("Removing tiller")
		Expect(helm.DeleteTiller(false)).NotTo(HaveOccurred())
	})

	It("Should be possible to create/delete/update and check status of wordpress chart", func() {
		By("Install chart stable/wordpress")
		releaseName, err := helm.Install("stable/wordpress")
		Expect(err).NotTo(HaveOccurred())
		By("Check status of release " + releaseName)
		Expect(helm.Status(releaseName)).NotTo(HaveOccurred())
	})
})
