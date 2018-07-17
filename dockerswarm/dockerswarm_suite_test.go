package dockerswarm_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestDockerswarm(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Dockerswarm Suite")
}
