package kubernetes

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/typed/apps/v1"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	"git.corp.adobe.com/abramowi/hyperion/lib/core"
)

type manager struct {
	clientset         *kubernetes.Clientset
	deploymentsClient v1.DeploymentInterface
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

func (mgr *manager) DeployApp(app core.App) (core.Operation, error) {
	k8sDeployment := getK8sDeployment(app)
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

func getK8sDeployment(app core.App) *appsv1.Deployment {
	countInt32 := int32(app.Count)
	k8sDeployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: app.ID},
		Spec: appsv1.DeploymentSpec{
			Replicas: &countInt32,
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{"appID": app.ID}},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"appID": app.ID}},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  app.ID,
							Image: app.Image,
							Ports: []apiv1.ContainerPort{
								{Name: "http", Protocol: apiv1.ProtocolTCP, ContainerPort: 80},
							},
						},
					},
				},
			},
		},
	}
	return k8sDeployment
}
