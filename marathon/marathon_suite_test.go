package marathon_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestMarathon(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Marathon Suite")
}
