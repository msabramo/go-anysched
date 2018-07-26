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
		writeJSONResponseFromFile(w, jsonResponseFilePath)
	}))
}

func writeJSONResponseFromFile(w http.ResponseWriter, jsonResponseFilePath string) {
	bytes, err := ioutil.ReadFile(jsonResponseFilePath)
	if err != nil {
		panic(err)
	}
	writeJSONResponseBytes(w, bytes)
}

func writeJSONResponseBytes(w http.ResponseWriter, bytes []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", len(bytes)))
	w.Write(bytes)
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

	Describe("Tasks", func() {
		var (
			manager hyperion.Manager
			ts      *httptest.Server
		)

		BeforeEach(func() {
			ts = NewTestServerJSONResponse("testdata/pods_list.json")
			manager = NewManagerWithTestServer(ts)
		})

		AfterEach(func() {
			ts.Close()
		})

		It("works", func() {
			tasks, err := manager.Tasks()
			Expect(err).ToNot(HaveOccurred())
			Expect(tasks).ToNot(BeNil())
			Expect(tasks).To(HaveLen(3))

			Expect(tasks[0].Name).To(Equal("httpbin-5d7c976bcd-9kjz5"))
			Expect(tasks[0].HostIP).To(Equal("10.0.2.15"))
			Expect(tasks[0].TaskIP).To(Equal("172.17.0.4"))
			Expect((*tasks[0].ReadyTime).Format(time.RFC3339)).To(Equal("2018-07-20T11:38:05-07:00"))

			Expect(tasks[1].Name).To(Equal("httpbin-5d7c976bcd-wmsvl"))
			Expect(tasks[1].HostIP).To(Equal("10.0.2.15"))
			Expect(tasks[1].TaskIP).To(Equal("172.17.0.5"))
			Expect((*tasks[1].ReadyTime).Format(time.RFC3339)).To(Equal("2018-07-20T11:38:07-07:00"))

			Expect(tasks[2].Name).To(Equal("httpbin-5d7c976bcd-xn6dx"))
			Expect(tasks[2].HostIP).To(Equal("10.0.2.15"))
			Expect(tasks[2].TaskIP).To(Equal("172.17.0.6"))
			Expect((*tasks[2].ReadyTime).Format(time.RFC3339)).To(Equal("2018-07-20T11:38:09-07:00"))
		})
	})

	Describe("SvcTasks", func() {
		var (
			manager hyperion.Manager
			ts      *httptest.Server
		)

		BeforeEach(func() {
			ts = NewTestServerJSONResponse("testdata/pods_list.json")
			manager = NewManagerWithTestServer(ts)
		})

		AfterEach(func() {
			ts.Close()
		})

		It("works", func() {
			tasks, err := manager.SvcTasks(hyperion.SvcCfg{ID: "httpbin"})
			Expect(err).ToNot(HaveOccurred())
			Expect(tasks).ToNot(BeNil())
			Expect(tasks).To(HaveLen(3))

			Expect(tasks[0].Name).To(Equal("httpbin-5d7c976bcd-9kjz5"))
			Expect(tasks[0].HostIP).To(Equal("10.0.2.15"))
			Expect(tasks[0].TaskIP).To(Equal("172.17.0.4"))
			Expect((*tasks[0].ReadyTime).Format(time.RFC3339)).To(Equal("2018-07-20T11:38:05-07:00"))

			Expect(tasks[1].Name).To(Equal("httpbin-5d7c976bcd-wmsvl"))
			Expect(tasks[1].HostIP).To(Equal("10.0.2.15"))
			Expect(tasks[1].TaskIP).To(Equal("172.17.0.5"))
			Expect((*tasks[1].ReadyTime).Format(time.RFC3339)).To(Equal("2018-07-20T11:38:07-07:00"))

			Expect(tasks[2].Name).To(Equal("httpbin-5d7c976bcd-xn6dx"))
			Expect(tasks[2].HostIP).To(Equal("10.0.2.15"))
			Expect(tasks[2].TaskIP).To(Equal("172.17.0.6"))
			Expect((*tasks[2].ReadyTime).Format(time.RFC3339)).To(Equal("2018-07-20T11:38:09-07:00"))
		})
	})

	Describe("DeploySvc", func() {
		var (
			manager hyperion.Manager
			ts      *httptest.Server
		)

		BeforeEach(func() {
			ts = NewTestServerJSONResponse("testdata/deployment_create.json")
			manager = NewManagerWithTestServer(ts)
		})

		AfterEach(func() {
			ts.Close()
		})

		It("works", func() {
			svcCfg := hyperion.SvcCfg{ID: "httpbin", Image: "citizenstig/httpbin", Count: 3}
			deployment, err := manager.DeploySvc(svcCfg)
			Expect(err).ToNot(HaveOccurred())
			Expect(deployment).ToNot(BeNil())
		})
	})
})
