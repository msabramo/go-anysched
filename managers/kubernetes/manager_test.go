package kubernetes

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("kubernetes/manager.go", func() {
	Describe("NewManager", func() {
		It("works with a valid URL", func() {
			manager, err := NewManager("http://1.2.3.4:8080")
			Expect(err).ToNot(HaveOccurred())
			Expect(manager).ToNot(BeNil())
		})

		It("works if URL is blank but KUBECONFIG is set", func() {
			oldKubeConfig := os.Getenv("KUBECONFIG")
			os.Setenv("KUBECONFIG", "../../etc/kubeconfigs/minikube.kubeconfig")
			defer func() { os.Setenv("KUBECONFIG", oldKubeConfig) }()
			manager, err := NewManager("")
			Expect(err).ToNot(HaveOccurred())
			Expect(manager).ToNot(BeNil())
		})

		It("fails with an invalid URL", func() {
			manager, err := NewManager(":::::---!@#$%")
			Expect(err).To(HaveOccurred())
			Expect(manager).To(BeNil())
		})

		It("fails with an invalid kubeconfig file", func() {
			manager, err := NewManager("/dev/does-not-exist")
			Expect(err).To(HaveOccurred())
			Expect(manager).To(BeNil())
		})
	})
})
