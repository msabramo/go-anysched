package marathon

import (
	"context"
	"fmt"
	"time"
)

type marathonDeployment struct {
	appID           string
	deploymentIDs   []string
	marathonManager marathonManager
}

func (d *marathonDeployment) Wait(ctx context.Context, timeout time.Duration) error {
	fmt.Printf("Wait() called with d = %+v\n", d)
	for _, deploymentID := range d.deploymentIDs {
		err := d.marathonManager.goMarathonClient.WaitOnDeployment(deploymentID, timeout)
		if err != nil {
			return err
		}
	}
	return nil
}
