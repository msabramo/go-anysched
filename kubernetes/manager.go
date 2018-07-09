package kubernetes

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"git.corp.adobe.com/abramowi/hyperion/core"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

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

func (k *k8sManager) DeployApp(app core.App) (core.Operation, error) {
	deploymentsClient := k.clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	countInt32 := int32(app.Count)
	deployment := &appsv1.Deployment{
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

func (k *k8sManager) DestroyApp(appID string) (core.Operation, error) {
	deploymentsClient := k.clientset.AppsV1().Deployments(apiv1.NamespaceDefault)

	err := deploymentsClient.Delete(appID, nil)
	return nil, err
}
