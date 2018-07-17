package kubernetes

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("kubernetes/manager.go", func() {
	Describe("NewManager", func() {
		It("works", func() {
			manager, err := NewManager("http://1.2.3.4:8080")
			Expect(err).ToNot(HaveOccurred())
			Expect(manager).ToNot(BeNil())
		})
	})
})
