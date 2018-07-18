package nomad_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestNomad(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nomad Suite")
}
