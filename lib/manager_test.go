package hyperion

import (
	"context"
	"time"

	"git.corp.adobe.com/abramowi/hyperion/lib/core"
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
			app := core.App{
				ID:    "my-app",
				Image: "citizenstig/httpbin:latest",
				Count: 2,
			}
			operation, err := manager.DeployApp(app)
			Expect(err).ToNot(HaveOccurred())
			Expect(operation).ToNot(BeNil())

			if asyncOperation, ok := operation.(core.AsyncOperation); ok && asyncOperation != nil {
				asyncOperation.Wait(ctx, 60*time.Second)
			}

			time.Sleep(10 * time.Second)

			destroyOperation, err := manager.DestroyApp("my-app")
			Expect(err).ToNot(HaveOccurred())
			Expect(destroyOperation).ToNot(BeNil())
		})
	})
})
