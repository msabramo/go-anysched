package kubernetes

import (
	"fmt"
	"time"

	appsv1 "k8s.io/api/apps/v1"
)

type deployment struct {
	appsv1Deployment *appsv1.Deployment
}

func (dep deployment) String() string {
	return fmt.Sprintf(
		"<kubernetes.deployment name=%q uid=%q creationTimestamp=%q />",
		dep.Name(), dep.UID(), dep.CreationTimestamp(),
	)
}

func (dep deployment) Name() string {
	return dep.appsv1Deployment.GetObjectMeta().GetName()
}

func (dep deployment) UID() string {
	return string(dep.appsv1Deployment.GetObjectMeta().GetUID())
}

func (dep deployment) CreationTimestamp() time.Time {
	return dep.appsv1Deployment.GetObjectMeta().GetCreationTimestamp().Time
}
