package nomad

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("nomad/manager.go", func() {
	Describe("NewManager", func() {
		It("works", func() {
			manager, err := NewManager("http://1.2.3.4:8080")
			Expect(err).ToNot(HaveOccurred())
			Expect(manager).ToNot(BeNil())
		})

		It("returns an error if address is invalid", func() {
			manager, err := NewManager(":")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("api.NewClient failed"))
			Expect(manager).To(BeNil())
		})
	})
})
