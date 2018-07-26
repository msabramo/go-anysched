package kubernetes

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"git.corp.adobe.com/abramowi/hyperion"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func NewTestServerJSONResponse(jsonResponseFilePath string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := ioutil.ReadFile(jsonResponseFilePath)
		if err != nil {
			panic(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Length", fmt.Sprintf("%d", len(body)))
		w.WriteHeader(200)
		w.Write(body)
	}))
}

func NewManagerWithTestServer(ts *httptest.Server) hyperion.Manager {
	manager, err := NewManager(ts.URL)
	if err != nil {
		panic(err)
	}
	return manager
}

var _ = Describe("kubernetes/manager.go", func() {
	Describe("NewManager", func() {
		It("works with a valid URL", func() {
			manager, err := NewManager("http://1.2.3.4:8080")
			Expect(err).ToNot(HaveOccurred())
			Expect(manager).ToNot(BeNil())
		})

		It("works if URL is blank but KUBECONFIG is set", func() {
			oldKubeConfig := os.Getenv("KUBECONFIG")
			os.Setenv("KUBECONFIG", "../../etc/kubeconfigs/minikube.kubeconfig")
			defer func() { os.Setenv("KUBECONFIG", oldKubeConfig) }()
			manager, err := NewManager("")
			Expect(err).ToNot(HaveOccurred())
			Expect(manager).ToNot(BeNil())
		})

		It("fails with an invalid URL", func() {
			manager, err := NewManager(":::::---!@#$%")
			Expect(err).To(HaveOccurred())
			Expect(manager).To(BeNil())
		})

		It("fails with a non-existent kubeconfig file", func() {
			manager, err := NewManager("/dev/does-not-exist")
			Expect(err).To(HaveOccurred())
			Expect(manager).To(BeNil())
		})

		It("fails with KUBECONFIG set to a non-existent kubeconfig file", func() {
			oldKubeConfig := os.Getenv("KUBECONFIG")
			os.Setenv("KUBECONFIG", "/dev/does-not-exist")
			defer func() { os.Setenv("KUBECONFIG", oldKubeConfig) }()
			manager, err := NewManager("")
			Expect(err).To(HaveOccurred())
			Expect(manager).To(BeNil())
		})

		It("fails with a garbage kubeconfig file 1", func() {
			manager, err := NewManager("/dev/null")
			Expect(err).To(HaveOccurred())
			Expect(manager).To(BeNil())
		})

		It("fails with KUBECONFIG set to a garbage kubeconfig file 1", func() {
			oldKubeConfig := os.Getenv("KUBECONFIG")
			os.Setenv("KUBECONFIG", "/dev/null")
			defer func() { os.Setenv("KUBECONFIG", oldKubeConfig) }()
			manager, err := NewManager("")
			Expect(err).To(HaveOccurred())
			Expect(manager).To(BeNil())
		})

		It("fails with a garbage kubeconfig file 2", func() {
			manager, err := NewManager("/etc/passwd")
			Expect(err).To(HaveOccurred())
			Expect(manager).To(BeNil())
		})

		It("fails with KUBECONFIG set to a garbage kubeconfig file 2", func() {
			oldKubeConfig := os.Getenv("KUBECONFIG")
			os.Setenv("KUBECONFIG", "/etc/passwd")
			defer func() { os.Setenv("KUBECONFIG", oldKubeConfig) }()
			manager, err := NewManager("")
			Expect(err).To(HaveOccurred())
			Expect(manager).To(BeNil())
		})
	})

	Describe("Svcs", func() {
		var (
			manager hyperion.Manager
			ts      *httptest.Server
		)

		BeforeEach(func() {
			ts = NewTestServerJSONResponse("testdata/deployments_list.json")
			manager = NewManagerWithTestServer(ts)
		})

		AfterEach(func() {
			ts.Close()
		})

		It("works", func() {
			svcs, err := manager.Svcs()
			Expect(err).ToNot(HaveOccurred())
			Expect(svcs).ToNot(BeNil())
			Expect(svcs).To(HaveLen(1))
			Expect(svcs[0].ID).To(Equal("httpbin"))
			Expect(*svcs[0].TasksRunning).To(Equal(3))
			Expect(*svcs[0].TasksHealthy).To(Equal(3))
			Expect(*svcs[0].TasksUnhealthy).To(Equal(0))
			Expect((*svcs[0].CreationTime).Format(time.RFC3339)).To(Equal("2018-07-20T11:38:03-07:00"))
		})
	})
})
