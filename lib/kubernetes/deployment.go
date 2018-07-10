package kubernetes

import (
	appsv1 "k8s.io/api/apps/v1"
)

type k8sDeployment struct {
	appsv1Deployment *appsv1.Deployment
}
