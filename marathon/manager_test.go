package marathon

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("marathon/manager.go", func() {
	Describe("NewManager", func() {
		It("works", func() {
			manager, err := NewManager("http://1.2.3.4:8080")
			Expect(err).ToNot(HaveOccurred())
			Expect(manager).ToNot(BeNil())
		})

		It("returns an error if address is blank", func() {
			manager, err := NewManager("")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("goMarathon.NewClient failed"))
			Expect(manager).To(BeNil())
		})
	})
})
