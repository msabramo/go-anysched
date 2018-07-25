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

		It("fails with a non-existent kubeconfig file", func() {
			manager, err := NewManager("/dev/does-not-exist")
			Expect(err).To(HaveOccurred())
			Expect(manager).To(BeNil())
		})

		It("fails with KUBECONFIG set to a non-existent kubeconfig file", func() {
			oldKubeConfig := os.Getenv("KUBECONFIG")
			os.Setenv("KUBECONFIG", "/dev/does-not-exist")
			defer func() { os.Setenv("KUBECONFIG", oldKubeConfig) }()
			manager, err := NewManager("")
			Expect(err).To(HaveOccurred())
			Expect(manager).To(BeNil())
		})

		It("fails with a garbage kubeconfig file 1", func() {
			manager, err := NewManager("/dev/null")
			Expect(err).To(HaveOccurred())
			Expect(manager).To(BeNil())
		})

		It("fails with KUBECONFIG set to a garbage kubeconfig file 1", func() {
			oldKubeConfig := os.Getenv("KUBECONFIG")
			os.Setenv("KUBECONFIG", "/dev/null")
			defer func() { os.Setenv("KUBECONFIG", oldKubeConfig) }()
			manager, err := NewManager("")
			Expect(err).To(HaveOccurred())
			Expect(manager).To(BeNil())
		})

		It("fails with a garbage kubeconfig file 2", func() {
			manager, err := NewManager("/etc/passwd")
			Expect(err).To(HaveOccurred())
			Expect(manager).To(BeNil())
		})

		It("fails with KUBECONFIG set to a garbage kubeconfig file 2", func() {
			oldKubeConfig := os.Getenv("KUBECONFIG")
			os.Setenv("KUBECONFIG", "/etc/passwd")
			defer func() { os.Setenv("KUBECONFIG", oldKubeConfig) }()
			manager, err := NewManager("")
			Expect(err).To(HaveOccurred())
			Expect(manager).To(BeNil())
		})
	})
})
