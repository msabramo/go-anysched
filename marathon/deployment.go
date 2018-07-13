package marathon

import (
	"context"
	"fmt"
	"time"

	goMarathon "github.com/gambol99/go-marathon"
	"github.com/pkg/errors"

	"git.corp.adobe.com/abramowi/hyperion/core"
)

type marathonDeploymentOperation struct {
	*manager
	appID           string
	deploymentIDs   []string
	timeoutDuration time.Duration
}

// GetProperties returns a map with all labels, annotations, and basic
// properties like name or uid
func (d *marathonDeploymentOperation) GetProperties() (propertiesMap map[string]interface{}) {
	return propertiesMap
}

func (d *marathonDeploymentOperation) GetStatus() (status *core.Status, err error) {
	opts := &goMarathon.GetAppOpts{
		Embed: []string{"app.tasks", "app.counts", "app.deployments", "app.readiness", "app.lastTaskFailure", "app.taskStats"},
	}
	app, err := d.manager.goMarathonClient.ApplicationBy(d.appID, opts)
	if err != nil {
		return nil, errors.Wrapf(err,
			"marathon.marathonDeploymentOperation.GetStatus: goMarathonClient.Application(%q) failed",
			d.appID)
	}
	// fmt.Printf("*** GetStatus: app = %+v\n", app)
	if !app.AllTaskRunning() {
		return &core.Status{
			ClientTime: time.Now(),
			// LastTransitionTime: lastTransitionTime,
			LastUpdateTime: time.Now(),
			Msg:            fmt.Sprintf("Not all tasks running. %d task(s) running.", app.TasksRunning),
			Done:           false,
		}, nil
	}
	return &core.Status{
		ClientTime: time.Now(),
		// LastTransitionTime: lastTransitionTime,
		LastUpdateTime: time.Now(),
		Msg:            fmt.Sprintf("All tasks running. %d task(s) running.", app.TasksRunning),
		Done:           true,
	}, nil
}

func (d *marathonDeploymentOperation) Wait(ctx context.Context) (result interface{}, err error) {
	fmt.Printf("Wait() called with d = %+v\n", d)
	for _, deploymentID := range d.deploymentIDs {
		err = d.manager.goMarathonClient.WaitOnDeployment(deploymentID, d.timeoutDuration)
		if err != nil {
			return nil, errors.Wrapf(err,
				"marathon.marathonDeploymentOperation.Wait: goMarathonClient.WaitOnDeployment(%q, %v) failed",
				deploymentID, d.timeoutDuration)
		}
	}
	return nil, nil
}
