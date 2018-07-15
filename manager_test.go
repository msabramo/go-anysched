package hyperion

import (
	"context"
	"time"

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
	})

	AfterEach(func() {
		mockCtrl.Finish()
	})

	Describe("deploying to Marathon", func() {
		It("deploys a service to Marathon as a Marathon application", func() {
			manager, err := NewManager(
				ManagerConfig{
					Type:    "marathon",
					Address: "http://127.0.0.1:8080",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			svc := SvcCfg{
				ID:    "my-svc",
				Image: "citizenstig/httpbin:latest",
				Count: 2,
			}
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
