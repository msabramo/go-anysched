package kubernetes

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/msabramo/go-anysched"
)

func NewTestServerJSONResponse(jsonResponseFilePath string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSONResponseFromFile(w, jsonResponseFilePath)
	}))
}

func NewTestServerJSONResponses(jsonResponseFilePaths ...string) *httptest.Server {
	count := 0
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSONResponseFromFile(w, jsonResponseFilePaths[count])
		count++
		if count >= len(jsonResponseFilePaths) {
			count = len(jsonResponseFilePaths) - 1
		}
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

func NewManagerWithTestServer(ts *httptest.Server) anysched.Manager {
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
			manager anysched.Manager
			ts      *httptest.Server
		)

		Context("healthy k8s", func() {
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

		Context("unhealthy k8s", func() {
			BeforeEach(func() {
				ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(500)
				}))
				manager = NewManagerWithTestServer(ts)
			})

			AfterEach(func() {
				ts.Close()
			})

			It("works", func() {
				svcs, err := manager.Svcs()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("deploymentsClient.List failed"))
				Expect(svcs).To(BeNil())
			})
		})
	})

	Describe("Tasks", func() {
		var (
			manager anysched.Manager
			ts      *httptest.Server
		)

		Context("healthy k8s", func() {
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

		Context("unhealthy k8s", func() {
			BeforeEach(func() {
				ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(500)
				}))
				manager = NewManagerWithTestServer(ts)
			})

			AfterEach(func() {
				ts.Close()
			})

			It("works", func() {
				tasks, err := manager.Tasks()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("podsClient.List failed"))
				Expect(tasks).To(BeNil())
			})
		})
	})

	Describe("SvcTasks", func() {
		var (
			manager anysched.Manager
			ts      *httptest.Server
		)

		Context("healthy k8s", func() {
			BeforeEach(func() {
				ts = NewTestServerJSONResponse("testdata/pods_list.json")
				manager = NewManagerWithTestServer(ts)
			})

			AfterEach(func() {
				ts.Close()
			})

			It("works", func() {
				tasks, err := manager.SvcTasks(anysched.SvcCfg{ID: "httpbin"})
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

		Context("unhealthy k8s", func() {
			BeforeEach(func() {
				ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(500)
				}))
				manager = NewManagerWithTestServer(ts)
			})

			AfterEach(func() {
				ts.Close()
			})

			It("works", func() {
				tasks, err := manager.SvcTasks(anysched.SvcCfg{ID: "httpbin"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("podsClient.List failed"))
				Expect(tasks).To(BeNil())
			})
		})
	})

	Describe("DeploySvc", func() {
		var (
			manager anysched.Manager
			ts      *httptest.Server
			svcCfg  anysched.SvcCfg
		)

		Context("successful deploy", func() {
			BeforeEach(func() {
				ts = NewTestServerJSONResponse("testdata/deployment_create.json")
				manager = NewManagerWithTestServer(ts)
				svcCfg = anysched.SvcCfg{ID: "httpbin", Image: "citizenstig/httpbin", Count: 3}
			})

			AfterEach(func() {
				ts.Close()
			})

			It("works", func() {
				deployment, err := manager.DeploySvc(svcCfg)
				Expect(err).ToNot(HaveOccurred())
				Expect(deployment).ToNot(BeNil())
			})
		})

		Context("k8s deployment creation fails with HTTP 500", func() {
			BeforeEach(func() {
				ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(500)
				}))
				manager = NewManagerWithTestServer(ts)
				svcCfg = anysched.SvcCfg{ID: "httpbin", Image: "citizenstig/httpbin", Count: 3}
			})

			AfterEach(func() {
				ts.Close()
			})

			It("works", func() {
				deployment, err := manager.DeploySvc(svcCfg)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("deploymentsClient.Create failed"))
				Expect(deployment).To(BeNil())
			})
		})
	})

	Describe("DestroySvc", func() {
		var (
			manager anysched.Manager
			ts      *httptest.Server
		)

		Context("successful destroy", func() {
			BeforeEach(func() {
				ts = NewTestServerJSONResponse("testdata/deployment_destroy_httpbin.json")
				manager = NewManagerWithTestServer(ts)
			})

			AfterEach(func() {
				ts.Close()
			})

			It("works", func() {
				destroy, err := manager.DestroySvc("httpbin")
				Expect(err).ToNot(HaveOccurred())
				Expect(destroy).To(BeNil())
			})
		})

		Context("k8s deployment destroy fails with HTTP 500", func() {
			BeforeEach(func() {
				ts = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(500)
				}))
				manager = NewManagerWithTestServer(ts)
			})

			AfterEach(func() {
				ts.Close()
			})

			It("works", func() {
				destroy, err := manager.DestroySvc("httpbin")
				Expect(err.Error()).To(ContainSubstring("deploymentsClient.Delete failed"))
				Expect(destroy).To(BeNil())
			})
		})
	})

	Context("a deployment exists", func() {
		var (
			ts           *httptest.Server
			myDeployment deployment
		)

		getDeployment := func() deployment {
			tsTmp := NewTestServerJSONResponse("testdata/deployment_create.json")
			defer tsTmp.Close()
			mgr := NewManagerWithTestServer(ts)
			svcCfg := anysched.SvcCfg{ID: "httpbin", Image: "citizenstig/httpbin", Count: 3}
			deploySvcResult, _ := mgr.DeploySvc(svcCfg)
			return deploySvcResult.(deployment)
		}

		BeforeEach(func() {
			ts = NewTestServerJSONResponse("testdata/deployment_get_httpbin.json")
			myDeployment = getDeployment()
		})

		AfterEach(func() {
			ts.Close()
		})

		Describe("isDone", func() {
			It("works", func() {
				Expect(myDeployment.isDone()).To(BeTrue())
			})
		})

		Describe("String", func() {
			It("works", func() {
				Expect(myDeployment.String()).To(Equal(`<kubernetes.deployment name="httpbin" ` +
					`uid="a4487ed1-9082-11e8-a0ad-080027aa669d" creationTimestamp="2018-07-25T20:19:01-07:00" />`))
			})
		})

		Describe("GetStatus", func() {
			It("works", func() {
				status, err := myDeployment.GetStatus()
				Expect(err).ToNot(HaveOccurred())
				Expect(status).ToNot(BeNil())
				Expect(status.Done).To(BeTrue())
				Expect(status.LastUpdateTime.Format(time.RFC3339)).To(Equal("2018-07-25T20:19:07-07:00"))
				Expect(status.LastTransitionTime.Format(time.RFC3339)).To(Equal("2018-07-25T20:19:07-07:00"))
				Expect(status.Msg).To(Equal(`Deployment "httpbin" successfully rolled out. ` +
					`3 of 3 updated replicas are available.`))
			})
		})

		Describe("GetProperties", func() {
			It("works", func() {
				props := myDeployment.GetProperties()
				Expect(props).ToNot(BeNil())
				Expect(props["name"]).To(Equal("httpbin"))
				Expect(props["uid"]).To(Equal(types.UID("a4487ed1-9082-11e8-a0ad-080027aa669d")))
				Expect(props["creationTimestamp"]).To(Equal("2018-07-25T20:19:01-07:00"))
				Expect(props["namespace"]).To(Equal("default"))
				Expect(props["generation"]).To(Equal(int64(1)))
				Expect(props["resourceVersion"]).To(Equal("222164"))
				Expect(props["annotations.deployment.kubernetes.io/revision"]).To(Equal("1"))
				Expect(props["selfLink"]).To(Equal("/apis/apps/v1/namespaces/default/deployments/httpbin"))
			})
		})
	})

	Context("a deployment that progresses", func() {
		var (
			ts           *httptest.Server
			myDeployment deployment
		)

		getDeployment := func() deployment {
			tsTmp := NewTestServerJSONResponse("testdata/deployment_create.json")
			defer tsTmp.Close()
			mgr := NewManagerWithTestServer(ts)
			svcCfg := anysched.SvcCfg{ID: "httpbin", Image: "citizenstig/httpbin", Count: 3}
			deploySvcResult, _ := mgr.DeploySvc(svcCfg)
			return deploySvcResult.(deployment)
		}

		BeforeEach(func() {
			ts = NewTestServerJSONResponses(
				"testdata/deployment_get_httpbin.json",
			)
			myDeployment = getDeployment()
		})

		AfterEach(func() {
			ts.Close()
		})

		Describe("Wait", func() {
			It("works", func() {
				_, err := myDeployment.Wait(context.Background())
				Expect(err).ToNot(HaveOccurred())
			})
		})
	})

	Context("a deployment that has replicas unavailable", func() {
		var (
			ts           *httptest.Server
			myDeployment deployment
		)

		getDeployment := func() deployment {
			tsTmp := NewTestServerJSONResponse("testdata/deployment_create.json")
			defer tsTmp.Close()
			mgr := NewManagerWithTestServer(ts)
			timeout := time.Duration(6 * time.Second)
			svcCfg := anysched.SvcCfg{
				ID:    "httpbin",
				Image: "citizenstig/httpbin",
				Count: 3,
				DeployTimeoutDuration: &timeout,
			}
			deploySvcResult, _ := mgr.DeploySvc(svcCfg)
			return deploySvcResult.(deployment)
		}

		BeforeEach(func() {
			ts = NewTestServerJSONResponses(
				"testdata/deployment_old_generation.json",
				"testdata/deployment_old_generation.json",
				"testdata/deployment_old_generation.json",
				"testdata/deployment_not_all_updated.json",
				"testdata/deployment_not_all_updated.json",
				"testdata/deployment_not_all_updated.json",
				"testdata/deployment_unavailable_replicas.json",
				"testdata/deployment_unavailable_replicas.json",
				"testdata/deployment_unavailable_replicas.json",
				"testdata/deployment_fail_not_progressing.json",
			)
			myDeployment = getDeployment()
		})

		AfterEach(func() {
			ts.Close()
		})

		Describe("GetStatus", func() {
			It("works", func() {
				var (
					status *anysched.OperationStatus
					err    error
				)

				for {
					status, err = myDeployment.GetStatus()
					if err != nil {
						break
					}
					Expect(status).ToNot(BeNil())
				}

				Expect(err).To(MatchError(`deployment "httpbin" exceeded its progress deadline`))
			})
		})

		Describe("Wait", func() {
			It("works", func() {
				_, err := myDeployment.Wait(context.Background())
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Timed out after 6s"))
			})
		})
	})
})
