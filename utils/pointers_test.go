package utils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("utils/pointers.go", func() {
	Describe("Sptr", func() {
		It("works on non-empty string", func() {
			Expect(*Sptr("foo")).To(Equal("foo"))
		})

		It("works on empty string", func() {
			Expect(*Sptr("")).To(Equal(""))
		})
	})

	Describe("Iptr", func() {
		It("works", func() {
			Expect(*Iptr(56)).To(Equal(56))
		})
	})

	Describe("Fptr", func() {
		It("works", func() {
			Expect(*Fptr(56.78)).To(Equal(56.78))
		})
	})

	Describe("Bptr", func() {
		It("works", func() {
			Expect(*Bptr(true)).To(Equal(true))
			Expect(*Bptr(false)).To(Equal(false))
		})
	})

	Describe("StringFromSptr", func() {
		It("works on non-empty string", func() {
			Expect(StringFromSptr(Sptr("foo"))).To(Equal("&`foo`"))
		})

		It("works on empty string", func() {
			Expect(StringFromSptr(Sptr(""))).To(Equal("&``"))
		})

		It("works on nil", func() {
			Expect(StringFromSptr(nil)).To(Equal("<nil>"))
		})
	})
})
