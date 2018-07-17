package utils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("utils/templates.go", func() {
	Describe("RenderTemplateToBytes", func() {
		It("works", func() {
			text := `The {{.Animal}} {{.Sound}}ed`
			data := struct{ Animal, Sound string }{Animal: "dog", Sound: "bark"}
			bytes, err := RenderTemplateToBytes("my-template", text, data)
			Expect(err).ToNot(HaveOccurred())
			Expect(string(bytes)).To(Equal("The dog barked"))
		})

		It("returns an error if the template cannot be parsed", func() {
			text := `The {{{{.Animal}} {{.Sound}}ed`
			data := struct{ Animal, Sound string }{Animal: "dog", Sound: "bark"}
			_, err := RenderTemplateToBytes("my-template", text, data)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("RenderGoTemplateToString: template.Parse failed"))
		})

		It("returns an error if the template cannot be executed", func() {
			text := `The {{.Animal}} {{.Sound}}ed`
			data := struct{ Animal, Sound func() string }{Animal: func() string { return "meow" }}
			_, err := RenderTemplateToBytes("my-template", text, data)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("RenderGoTemplateToString: template.Execute failed"))
		})
	})
})
