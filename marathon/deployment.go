package marathon

import (
	"context"
	"errors"
	"fmt"
	"time"

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
	return nil, errors.New("marathon.deployment.GetStatus: Not implemented yet")
}

func (d *marathonDeploymentOperation) Wait(ctx context.Context) (result interface{}, err error) {
	fmt.Printf("Wait() called with d = %+v\n", d)
	for _, deploymentID := range d.deploymentIDs {
		err = d.manager.goMarathonClient.WaitOnDeployment(deploymentID, d.timeoutDuration)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}
