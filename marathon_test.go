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
			appDeployer, err := NewAppDeployer(
				AppDeployerConfig{
					Type:    "marathon",
					Address: "http://127.0.0.1:8080",
				},
			)
			Expect(err).ToNot(HaveOccurred())
			marathonApp := appDeployer.NewApp().
				SetID("my-app").
				SetDockerImage("citizenstig/httpbin:latest").
				SetCount(2)
			operation, err := appDeployer.DeployApp(marathonApp)
			Expect(err).ToNot(HaveOccurred())
			Expect(operation).ToNot(BeNil())

			if asyncOperation, ok := operation.(AsyncOperation); ok && asyncOperation != nil {
				asyncOperation.Wait(ctx, 60*time.Second)
			}

			time.Sleep(10 * time.Second)

			destroyOperation, err := appDeployer.DestroyApp("my-app")
			Expect(err).ToNot(HaveOccurred())
			Expect(destroyOperation).ToNot(BeNil())
		})
	})
})
