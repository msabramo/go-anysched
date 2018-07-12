package hyperion

import (
	"context"
	"time"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("app.go", func() {
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

	Describe("MarathonManager.CreateApplication", func() {
		It("deploys an application to Marathon", func() {
			manager, err := NewManager(
				ManagerConfig{
					Type:    "marathon",
					Address: "http://127.0.0.1:8080",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			app := App{
				ID:    "my-app",
				Image: "citizenstig/httpbin:latest",
				Count: 2,
			}
			operation, err := manager.DeployApp(app)
			Expect(err).ToNot(HaveOccurred())
			Expect(operation).ToNot(BeNil())

			operation.Wait(ctx)

			time.Sleep(10 * time.Second)

			destroyOperation, err := manager.DestroyApp("my-app")
			Expect(err).ToNot(HaveOccurred())
			Expect(destroyOperation).ToNot(BeNil())
		})
	})
})
