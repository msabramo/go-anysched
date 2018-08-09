// +build integration_tests

package marathon

import (
	"context"
	"fmt"
	"time"

	"github.com/msabramo/go-anysched"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Marathon integration test", func() {
	var (
		mockCtrl *gomock.Controller
		ctx      = context.Background()
	)

	BeforeEach(func() {
		mockCtrl = gomock.NewController(GinkgoT())
		fmt.Printf("registered marathon\n")
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("deploying to Marathon", func() {
		It("deploys a service to Marathon as a Marathon application", func() {
			manager, err := anysched.NewManager(
				anysched.ManagerConfig{Type: "marathon", Address: "http://127.0.0.1:8080"},
			)
			Expect(err).ToNot(HaveOccurred())
			svc := anysched.SvcCfg{ID: "my-svc", Image: "citizenstig/httpbin:latest", Count: 2}
			deployOperation, err := manager.DeploySvc(svc)
			Expect(err).ToNot(HaveOccurred())
			Expect(deployOperation).ToNot(BeNil())

			_, err = deployOperation.Wait(ctx)
			Expect(err).ToNot(HaveOccurred())

			time.Sleep(10 * time.Second)

			destroyOperation, err := manager.DestroySvc("my-svc")
			Expect(err).ToNot(HaveOccurred())
			Expect(destroyOperation).ToNot(BeNil())
		})
	})
})
