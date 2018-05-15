package hyperion

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	marathon "github.com/gambol99/go-marathon"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type marathonDeployment struct {
	appID           string
	deploymentIDs   []string
	marathonManager marathonManager
}

type k8sDeployment struct {
	appsv1Deployment *appsv1.Deployment
}

func (d *marathonDeployment) Wait(ctx context.Context, timeout time.Duration) error {
	fmt.Printf("Wait() called with d = %+v\n", d)
	for _, deploymentID := range d.deploymentIDs {
		err := d.marathonManager.marathonClient.WaitOnDeployment(deploymentID, timeout)
		if err != nil {
			return err
		}
	}
	return nil
}

type AppDeployerConfig struct {
	Type    string // e.g.: "marathon", "kubernetes", etc.
	Address string // e.g.: "http://127.0.0.1:8080"
}

type AppDeployer interface {
	DeployApp(App) (Operation, error)
	DestroyApp(appID string) (Operation, error)
}

func NewAppDeployer(a AppDeployerConfig) (appDeployer AppDeployer, err error) {
	switch a.Type {
	case "marathon":
		return NewMarathonManager(a.Address)
	case "kubernetes":
		return NewK8sManager(a.Address)
	default:
		return nil, fmt.Errorf("Unknown type: %q", a.Type)
	}
}

type marathonManager struct {
	marathonClient marathon.Marathon
	url            string
}

func NewMarathonManager(url string) (*marathonManager, error) {
	config := marathon.NewDefaultConfig()
	config.URL = url
	client, err := marathon.NewClient(config)
	if err != nil {
		return nil, err
	}
	m := &marathonManager{marathonClient: client, url: url}
	return m, nil
}

type k8sManager struct {
	clientset *kubernetes.Clientset
}

func k8SConfigFromKubeConfig() *rest.Config {
	var kubeconfig *string
	if home := os.Getenv("HOME"); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}
	return config
}

func NewK8sManager(url string) (*k8sManager, error) {
	clientset, err := kubernetes.NewForConfig(k8SConfigFromKubeConfig())
	if err != nil {
		return nil, err
	}
	// fmt.Printf("*** clientset = %+v\n", clientset)
	// namespaces, err := clientset.CoreV1().Namespaces().List(metav1.ListOptions{})
	// fmt.Printf("*** namespaces = %+v\n", namespaces)
	// pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	// fmt.Printf("*** pods = %+v\n", pods)

	k8s := &k8sManager{clientset: clientset}
	return k8s, nil
}

func (k *k8sManager) DeployApp(app App) (Operation, error) {
	deploymentsClient := k.clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	countInt32 := int32(app.Count)
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{Name: app.ID},
		Spec: appsv1.DeploymentSpec{
			Replicas: &countInt32,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"appID": app.ID},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"appID": app.ID},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  app.ID,
							Image: app.Image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "http",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		panic(err)
	}
	operation := k8sDeployment{appsv1Deployment: result}
	fmt.Printf("Created deployment %+v (%q).\n", operation, result.GetObjectMeta().GetName())

	return operation, err
}

func (k *k8sManager) DestroyApp(appID string) (Operation, error) {
	deploymentsClient := k.clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	err := deploymentsClient.Delete(appID, nil)
	return nil, err
}

func (m *marathonManager) goMarathonApp(app App) (gomApp *marathon.Application) {
	gomApp = marathon.NewDockerApplication()
	gomApp.ID = app.ID
	gomApp.Container.Docker.Container(app.Image)
	gomApp.Count(app.Count)
	return gomApp
}

func (m *marathonManager) deploymentIDs(gomApp *marathon.Application) (deploymentIDs []string) {
	marathonDeploymentIDs := gomApp.DeploymentIDs()
	deploymentIDs = make([]string, len(marathonDeploymentIDs))
	for i, marathonDeploymentID := range marathonDeploymentIDs {
		deploymentIDs[i] = marathonDeploymentID.DeploymentID
	}
	return deploymentIDs
}

func (m *marathonManager) DeployApp(app App) (Operation, error) {
	gomApp, err := m.marathonClient.CreateApplication(m.goMarathonApp(app))
	if err != nil {
		return nil, err
	}
	op := &marathonDeployment{
		appID:           gomApp.ID,
		deploymentIDs:   m.deploymentIDs(gomApp),
		marathonManager: *m,
	}
	return op, err
}

func (m *marathonManager) DestroyApp(appID string) (Operation, error) {
	force := false
	marathonDeploymentID, err := m.marathonClient.DeleteApplication(appID, force)
	if err != nil {
		return nil, err
	}
	op := &marathonDeployment{appID: appID, deploymentIDs: []string{marathonDeploymentID.DeploymentID}, marathonManager: *m}
	return op, err
}
