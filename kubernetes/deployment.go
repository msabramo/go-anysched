package kubernetes

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"git.corp.adobe.com/abramowi/hyperion/core"
)

const (
	timedOutReason = "ProgressDeadlineExceeded"
)

var (
	getDeployTimeoutDuration = func(app core.App) time.Duration {
		if app.DeployTimeoutDuration == nil {
			return 60 * time.Second
		}
		return *app.DeployTimeoutDuration
	}
)

// deployment implements the hyperion.Operation interface
type deployment struct {
	*appsv1.Deployment
	manager *manager
	app     core.App
}

func (dep deployment) String() string {
	return fmt.Sprintf(
		"<kubernetes.deployment name=%q uid=%q creationTimestamp=%q />",
		dep.GetName(), dep.GetUID(), dep.GetCreationTimestamp().Format(time.RFC3339),
	)
}

// GetProperties returns a map with all labels, annotations, and basic
// properties like name or uid
func (dep deployment) GetProperties() (propertiesMap map[string]interface{}) {
	propertiesMap = map[string]interface{}{}
	for key, val := range dep.GetLabels() {
		propertiesMap["labels."+key] = val
	}
	for key, val := range dep.GetAnnotations() {
		propertiesMap["annotations."+key] = val
	}
	propertiesMap["name"] = dep.GetName()
	propertiesMap["uid"] = dep.GetUID()
	propertiesMap["creationTimestamp"] = dep.GetCreationTimestamp().Format(time.RFC3339)
	propertiesMap["namespace"] = dep.GetNamespace()
	propertiesMap["generation"] = dep.GetGeneration()
	propertiesMap["generateName"] = dep.GetGenerateName()
	propertiesMap["clusterName"] = dep.GetClusterName()
	propertiesMap["resourceVersion"] = dep.GetResourceVersion()
	propertiesMap["selfLink"] = dep.GetSelfLink()
	propertiesMap["spec.strategy"] = dep.Spec.Strategy
	return propertiesMap
}

func (dep deployment) GetStatus() (status *core.Status, err error) {
	k8sDeployment, err := dep.manager.deploymentsClient.Get(dep.GetName(), metav1.GetOptions{})
	if err != nil {
		return nil, errors.Wrap(err, "kubernetes.deployment.GetStatus: deploymentsClient.Get failed")
	}
	if k8sDeployment.Generation <= k8sDeployment.Status.ObservedGeneration {
		cond := getDeploymentCondition(k8sDeployment.Status, appsv1.DeploymentProgressing)
		if cond != nil && cond.Reason == timedOutReason {
			return nil, fmt.Errorf("deployment %q exceeded its progress deadline", dep.GetName())
		}
		if k8sDeployment.Spec.Replicas != nil && k8sDeployment.Status.UpdatedReplicas < *k8sDeployment.Spec.Replicas {
			return notAllReplicasUpdatedStatus(k8sDeployment), nil
		}
		if k8sDeployment.Status.Replicas > k8sDeployment.Status.UpdatedReplicas {
			return oldReplicasPendingTerminationStatus(k8sDeployment), nil
		}
		if k8sDeployment.Status.AvailableReplicas < k8sDeployment.Status.UpdatedReplicas {
			return notAllReplicasAvailableStatus(k8sDeployment), nil
		}
		return deploymentSuccessStatus(k8sDeployment), nil
	}
	return deploymentSpecUpdateNotObservedStatus(k8sDeployment), nil
}

// getDeploymentCondition returns the condition with the provided type.
func getDeploymentCondition(
	status appsv1.DeploymentStatus,
	condType appsv1.DeploymentConditionType,
) *appsv1.DeploymentCondition {
	for i := range status.Conditions {
		c := status.Conditions[i]
		if c.Type == condType {
			return &c
		}
	}
	return nil
}

// getPodCondition returns the condition with the provided type.
func getPodCondition(
	status apiv1.PodStatus,
	condType apiv1.PodConditionType,
) *apiv1.PodCondition {
	for i := range status.Conditions {
		c := status.Conditions[i]
		if c.Type == condType {
			return &c
		}
	}
	return nil
}

func mostRecentConditionTimes(status appsv1.DeploymentStatus) (lastTransitionTime, lastUpdateTime time.Time) {
	for i := range status.Conditions {
		c := status.Conditions[i]
		if c.LastTransitionTime.Time.After(lastTransitionTime) {
			lastTransitionTime = c.LastTransitionTime.Time
		}
		if c.LastUpdateTime.Time.After(lastUpdateTime) {
			lastUpdateTime = c.LastUpdateTime.Time
		}
	}
	return lastTransitionTime, lastUpdateTime
}

func waitingForDeploymentMsg(k8sDeployment *appsv1.Deployment) string {
	return fmt.Sprintf("Waiting for deployment %q to finish", k8sDeployment.GetName())
}

func deploymentSpecUpdateNotObservedStatus(k8sDeployment *appsv1.Deployment) *core.Status {
	msg := "Waiting for deployment spec update to be observed..."
	return notDoneStatus(k8sDeployment, msg)
}

func notAllReplicasUpdatedStatus(k8sDeployment *appsv1.Deployment) *core.Status {
	msg := fmt.Sprintf("%d out of %d new replicas have been updated...",
		k8sDeployment.Status.UpdatedReplicas, *k8sDeployment.Spec.Replicas)
	return notDoneStatus(k8sDeployment, msg)
}

func notAllReplicasAvailableStatus(k8sDeployment *appsv1.Deployment) *core.Status {
	msg := fmt.Sprintf("%d of %d updated replicas are available...",
		k8sDeployment.Status.AvailableReplicas, k8sDeployment.Status.UpdatedReplicas)
	return notDoneStatus(k8sDeployment, msg)
}

func oldReplicasPendingTerminationStatus(k8sDeployment *appsv1.Deployment) *core.Status {
	msg := fmt.Sprintf("%d old replicas are pending termination...",
		k8sDeployment.Status.Replicas-k8sDeployment.Status.UpdatedReplicas)
	return notDoneStatus(k8sDeployment, msg)
}

func deploymentSuccessStatus(k8sDeployment *appsv1.Deployment) *core.Status {
	msg := fmt.Sprintf("Deployment %q successfully rolled out. %d of %d updated replicas are available.",
		k8sDeployment.GetName(), k8sDeployment.Status.AvailableReplicas, k8sDeployment.Status.UpdatedReplicas)
	return doneStatus(k8sDeployment, msg)
}

func notDoneStatus(k8sDeployment *appsv1.Deployment, msg string) *core.Status {
	msg = fmt.Sprintf("%s: %s", waitingForDeploymentMsg(k8sDeployment), msg)
	return status(k8sDeployment, msg, false)
}

func doneStatus(k8sDeployment *appsv1.Deployment, msg string) *core.Status {
	return status(k8sDeployment, msg, true)
}

func status(k8sDeployment *appsv1.Deployment, msg string, done bool) *core.Status {
	lastTransitionTime, lastUpdateTime := mostRecentConditionTimes(k8sDeployment.Status)
	return &core.Status{
		ClientTime:         time.Now(),
		LastTransitionTime: lastTransitionTime,
		LastUpdateTime:     lastUpdateTime,
		Msg:                msg,
		Done:               done,
	}
}

func (dep deployment) Wait(ctx context.Context) (result interface{}, err error) {
	timeout := getDeployTimeoutDuration(dep.app)
	ctx, _ = context.WithTimeout(ctx, timeout)

	for {
		select {
		case <-ctx.Done():
			return nil, errors.Wrapf(ctx.Err(), "kubernetes.deployment.Wait: Timed out after %s", timeout)
		case <-time.After(2 * time.Second):
			k8sDeployment, err := dep.manager.deploymentsClient.Get(dep.GetName(), metav1.GetOptions{})
			if err != nil {
				return nil, errors.Wrap(err, "kubernetes.deployment.Wait: deploymentsClient.Get failed")
			}
			fmt.Printf("*** k8sDeployment.Status = %+v\n", k8sDeployment.Status)
			if k8sDeployment.Status.ObservedGeneration == k8sDeployment.Generation {
				if k8sDeployment.Status.AvailableReplicas == k8sDeployment.Status.UpdatedReplicas {
					return deployment{manager: dep.manager, Deployment: k8sDeployment, app: dep.app}, nil
				}
			}
		}
	}
}
