package kubernetes

import (
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
)

type k8sDeployment struct {
	appsv1Deployment *appsv1.Deployment
}

func (self k8sDeployment) String() string {
	name := self.appsv1Deployment.GetObjectMeta().GetName()
	uid := self.appsv1Deployment.GetObjectMeta().GetUID()
	return fmt.Sprintf("<k8sDeployment name=%q uid=%q />", name, uid)
}
