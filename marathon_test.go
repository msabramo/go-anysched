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
			marathonManager, err := NewMarathonManager()
			Expect(err).ToNot(HaveOccurred())
			marathonApp := marathonManager.NewApp().
				SetID("my-app").
				SetDockerImage("citizenstig/httpbin:latest").
				SetCount(4)
			marathonApp, err = marathonManager.CreateApplication(marathonApp)
			Expect(err).ToNot(HaveOccurred())
			Expect(ctx).ToNot(BeNil())
			time.Sleep(10 * time.Second)
			deleteRequest := marathonApp.NewDeleteRequest()
			err = marathonManager.DeleteApplication(deleteRequest)
			Expect(err).ToNot(HaveOccurred())
		})
	})
})
