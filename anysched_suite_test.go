package anysched_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestAnySched.(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "AnySched. Suite")
}
