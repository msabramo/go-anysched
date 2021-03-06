package kubernetes

import (
	"os"
	"path/filepath"
	"sort"

	"github.com/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	tappsv1 "k8s.io/client-go/kubernetes/typed/apps/v1"
	tcorev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/msabramo/go-anysched"
	"github.com/msabramo/go-anysched/utils"
)

type manager struct {
	clientset         *kubernetes.Clientset
	deploymentsClient tappsv1.DeploymentInterface
	podsClient        tcorev1.PodInterface
	namespacesClient  tcorev1.NamespaceInterface
}

func init() {
	anysched.RegisterManagerType("kubernetes", NewManager)
}

// NewManager returns a Manager for Kubernetes.
func NewManager(url string) (anysched.Manager, error) {
	var (
		restConfig *rest.Config
		err        error
	)
	if url == "" || url == "kubeconfig" {
		restConfig, err = configFromKubeconfig()
	} else {
		restConfig, err = configFromURL(url)
	}
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes.NewManager: kubernetes.NewForConfig failed")
	}

	mgr := &manager{
		clientset:         clientset,
		deploymentsClient: clientset.AppsV1().Deployments(apiv1.NamespaceDefault),
		namespacesClient:  clientset.CoreV1().Namespaces(),
		podsClient:        clientset.CoreV1().Pods(apiv1.NamespaceDefault),
	}
	return mgr, nil
}

func configFromKubeconfig() (*rest.Config, error) {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", getKubeconfig())
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes.configFromKubeconfig: clientcmd.BuildConfigFromFlags failed")
	}
	return config, nil
}

func configFromURL(url string) (*rest.Config, error) {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags(url, getKubeconfig())
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes.configFromURL: clientcmd.BuildConfigFromFlags failed")
	}
	return config, nil
}

func getKubeconfig() string {
	if os.Getenv("KUBECONFIG") != "" {
		return os.Getenv("KUBECONFIG")
	}
	return filepath.Join(os.Getenv("HOME"), ".kube", "config")
}

// Svcs returns info about all running services
func (mgr *manager) Svcs() ([]anysched.Svc, error) {
	k8sDeploymentList, err := mgr.deploymentsClient.List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes.manager.Svcs: deploymentsClient.List failed")
	}
	svcs := make([]anysched.Svc, len(k8sDeploymentList.Items))
	for i := range k8sDeploymentList.Items {
		k8sDeployment := k8sDeploymentList.Items[i]
		tasksRunning := int(k8sDeployment.Status.Replicas)
		tasksHealthy := int(k8sDeployment.Status.AvailableReplicas)
		tasksUnhealthy := int(k8sDeployment.Status.UnavailableReplicas)
		creationTimestamp := k8sDeployment.GetCreationTimestamp().Time
		svcs[i] = anysched.Svc{
			ID:             k8sDeployment.GetName(),
			TasksRunning:   &tasksRunning,
			TasksHealthy:   &tasksHealthy,
			TasksUnhealthy: &tasksUnhealthy,
			CreationTime:   &creationTimestamp,
		}
	}
	return svcs, nil
}

// Tasks returns info about all running tasks
func (mgr *manager) Tasks() ([]anysched.Task, error) {
	k8sPodList, err := mgr.podsClient.List(metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes.manager.Tasks: podsClient.List failed")
	}
	tasks := make([]anysched.Task, len(k8sPodList.Items))
	for i, k8sPod := range k8sPodList.Items {
		tasks[i] = *taskFromK8SPod(k8sPod)
	}
	return tasks, nil
}

// SvcTasks returns info about the running tasks for a service
func (mgr *manager) SvcTasks(svcCfg anysched.SvcCfg) ([]anysched.Task, error) {
	k8sPodList, err := mgr.podsClient.List(metav1.ListOptions{LabelSelector: "appID=" + svcCfg.ID})
	if err != nil {
		return nil, errors.Wrapf(err, "kubernetes.manager.SvcTasks: podsClient.List failed for svcCfg.ID = %q", svcCfg.ID)
	}
	tasks := make([]anysched.Task, len(k8sPodList.Items))
	for i, k8sPod := range k8sPodList.Items {
		tasks[i] = *taskFromK8SPod(k8sPod)
	}
	sortTasksByReadyTime(tasks)
	return tasks, nil
}

func taskFromK8SPod(k8sPod apiv1.Pod) *anysched.Task {
	cond := getPodCondition(k8sPod.Status, apiv1.PodReady)
	if cond == nil {
		return nil
	}
	return &anysched.Task{
		Name:      k8sPod.GetName(),
		HostIP:    k8sPod.Status.HostIP,
		TaskIP:    k8sPod.Status.PodIP,
		ReadyTime: &cond.LastTransitionTime.Time,
	}
}

func sortTasksByReadyTime(tasks []anysched.Task) {
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].ReadyTime != nil &&
			tasks[j].ReadyTime != nil &&
			tasks[i].ReadyTime.Before(*tasks[j].ReadyTime)
	})
}

// DeploySvc takes a SvcCfg and deploys it, returning an Operation.
func (mgr *manager) DeploySvc(svcCfg anysched.SvcCfg) (anysched.Operation, error) {
	k8sDeploymentRequest, err := getK8sDeploymentRequest(svcCfg)
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes.manager.DeploySvc: getK8sDeploymentRequest failed")
	}
	k8sDeployment, err := mgr.deploymentsClient.Create(k8sDeploymentRequest)
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes.manager.DeploySvc: deploymentsClient.Create failed")
	}

	return deployment{manager: mgr, Deployment: k8sDeployment, svcCfg: svcCfg}, nil
}

// DestroySvc destroys a service.
func (mgr *manager) DestroySvc(svcID string) (anysched.Operation, error) {
	err := mgr.deploymentsClient.Delete(svcID, &metav1.DeleteOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes.manager.DestroySvc: deploymentsClient.Delete failed")
	}
	return nil, nil
}

func getK8sDeploymentRequest(svcCfg anysched.SvcCfg) (*appsv1.Deployment, error) {
	var k8sDeploymentRequest appsv1.Deployment
	data, err := utils.RenderTemplateToBytes("kubernetes-deployment", deploymentYAMLTemplateString, svcCfg)
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes.getK8sDeploymentRequest: RenderTemplateToBytes failed")
	}
	err = decodeYAMLOrJSON(data, &k8sDeploymentRequest)
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes.getK8sDeploymentRequest: decodeYAMLOrJSON failed")
	}
	return &k8sDeploymentRequest, nil
}

// decodeYAMLOrJSON takes as input `inYAMLOrJSONBytes`: a []byte with YAML or
// JSON and decodes into the parameter called `out`.
func decodeYAMLOrJSON(inYAMLOrJSONBytes []byte, out runtime.Object) error {
	var defaults *schema.GroupVersionKind
	_, _, err := scheme.Codecs.UniversalDeserializer().Decode(inYAMLOrJSONBytes, defaults, out)
	if err != nil {
		return errors.Wrap(err, "kubernetes.decodeYAMLOrJSON: UniversalDeserializer().Decode failed")
	}
	return nil
}

var deploymentYAMLTemplateString = `
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{.ID}}
spec:
  replicas: {{.Count}}
  selector:
    matchLabels:
      appID: {{.ID}}
  template:
    metadata:
      labels:
        appID: {{.ID}}
    spec:
      containers:
        - name: {{.ID}}
          image: {{.Image}}`
