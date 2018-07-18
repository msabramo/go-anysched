package marathon

import (
	"context"
	"fmt"
	"time"

	goMarathon "github.com/gambol99/go-marathon"
	"github.com/pkg/errors"

	"git.corp.adobe.com/abramowi/hyperion"
)

type deployment struct {
	*manager
	svcID                 string
	marathonDeploymentIDs []string
	timeoutDuration       time.Duration
}

// GetProperties returns a map with all labels, annotations, and basic
// properties like name or uid
func (d *deployment) GetProperties() (propertiesMap map[string]interface{}) {
	propertiesMap = map[string]interface{}{}
	if len(d.marathonDeploymentIDs) == 1 {
		marathonDeploymentID := d.marathonDeploymentIDs[0]
		propertiesMap["marathonDeploymentID"] = marathonDeploymentID
	}
	return propertiesMap
}

func (d *deployment) GetStatus() (status *hyperion.OperationStatus, err error) {
	opts := &goMarathon.GetAppOpts{
		Embed: []string{
			"app.tasks", "app.counts", "app.deployments",
			"app.readiness", "app.lastTaskFailure", "app.taskStats",
		},
	}
	goMarathonApp, err := d.manager.goMarathonClient.ApplicationBy(d.svcID, opts)
	if err != nil {
		return nil, errors.Wrapf(err,
			"marathon.deployment.GetStatus: goMarathonClient.ApplicationBy(%q) failed", d.svcID)
	}

	if !goMarathonApp.AllTaskRunning() {
		return notAllTasksRunningStatus(goMarathonApp), nil
	}
	return allTasksRunningStatus(goMarathonApp), nil
}

func notAllTasksRunningStatus(goMarathonApp *goMarathon.Application) *hyperion.OperationStatus {
	return statusWithTimestamps(
		&hyperion.OperationStatus{
			Msg:  fmt.Sprintf("Not all tasks running. %d task(s) running.", goMarathonApp.TasksRunning),
			Done: false,
		},
	)
}

func allTasksRunningStatus(goMarathonApp *goMarathon.Application) *hyperion.OperationStatus {
	return statusWithTimestamps(
		&hyperion.OperationStatus{
			Msg:  fmt.Sprintf("All tasks running. %d task(s) running.", goMarathonApp.TasksRunning),
			Done: true,
		},
	)
}

func statusWithTimestamps(status *hyperion.OperationStatus) *hyperion.OperationStatus {
	status.ClientTime = time.Now()
	status.LastUpdateTime = time.Now()
	// status.LastTransitionTime = lastTransitionTime
	return status
}

func (d *deployment) Wait(ctx context.Context) (result interface{}, err error) {
	for _, marathonDeploymentID := range d.marathonDeploymentIDs {
		err = d.manager.goMarathonClient.WaitOnDeployment(marathonDeploymentID, d.timeoutDuration)
		if err != nil {
			return nil, errors.Wrapf(err,
				"marathon.marathonDeploymentOperation.Wait: goMarathonClient.WaitOnDeployment(%q, %v) failed",
				marathonDeploymentID, d.timeoutDuration)
		}
	}
	return nil, nil
}
