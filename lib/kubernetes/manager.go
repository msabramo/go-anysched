package kubernetes

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/kubernetes/typed/apps/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"git.corp.adobe.com/abramowi/hyperion/lib/core"
	"git.corp.adobe.com/abramowi/hyperion/lib/utils"
)

type manager struct {
	clientset         *kubernetes.Clientset
	deploymentsClient v1.DeploymentInterface
}

func NewManager(url string) (*manager, error) {
	restConfig, err := configFromKubeconfig()
	if err != nil {
		return nil, err
	}
	clientset, err := kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes.NewManager: kubernetes.NewForConfig failed")
	}
	// fmt.Printf("*** clientset = %+v\n", clientset)
	// namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	// fmt.Printf("*** namespaces = %+v\n", namespaces)
	// pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	// fmt.Printf("*** pods = %+v\n", pods)

	mgr := &manager{
		clientset:         clientset,
		deploymentsClient: clientset.AppsV1().Deployments(apiv1.NamespaceDefault),
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

func getKubeconfig() string {
	var kubeconfig *string

	defaultKubeconfigFilePath := getDefaultKubeconfigFilePath(os.Getenv("HOME"))
	if defaultKubeconfigFilePath != "" {
		kubeconfig = flag.String("kubeconfig", defaultKubeconfigFilePath, "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	flag.Parse()

	return *kubeconfig
}

func getDefaultKubeconfigFilePath(homeDirPath string) string {
	if homeDirPath != "" {
		return filepath.Join(homeDirPath, ".kube", "config")
	}
	return ""
}

func (mgr *manager) DeployApp(app core.App) (core.Operation, error) {
	k8sDeployment, err := getK8sDeployment(app)
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes.manager.DeployApp: getK8sDeployment failed")
	}
	result, err := mgr.deploymentsClient.Create(k8sDeployment)
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes.manager.DeployApp: deploymentsClient.Create failed")
	}

	return deployment{appsv1Deployment: result}, nil
}

func (mgr *manager) DestroyApp(appID string) (core.Operation, error) {
	err := mgr.deploymentsClient.Delete(appID, nil)
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes.manager.DestroyApp: deploymentsClient.Delete failed")
	}
	return nil, nil
}

func getK8sDeployment(app core.App) (*appsv1.Deployment, error) {
	var k8sDeployment appsv1.Deployment
	data, err := utils.RenderTemplateToBytes("kubernetes-deployment", deploymentYAMLTemplateString, app)
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes.getK8sDeployment: RenderTemplateToBytes failed")
	}
	err = decodeYAMLOrJSON(data, &k8sDeployment)
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes.getK8sDeployment: decodeYAMLOrJSON failed")
	}
	return &k8sDeployment, nil
}

func decodeYAMLOrJSON(data []byte, into runtime.Object) error {
	var defaults *schema.GroupVersionKind
	_, _, err := scheme.Codecs.UniversalDeserializer().Decode(data, defaults, into)
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
