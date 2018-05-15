package hyperion

import (
	"context"
	"errors"
	"fmt"
	"time"

	marathon "github.com/gambol99/go-marathon"
)

type App interface {
	SetCount(int) App
	SetDockerImage(string) App
	SetID(string) App
}

type Operation interface{}

type AsyncOperation interface {
	Wait(ctx context.Context, timeout time.Duration) error // Wait for async operation to finish and return error or nil
}

type marathonDeployment struct {
	appID           string
	deploymentIDs   []string
	marathonManager marathonManager
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
	NewApp() App
	DeployApp(App) (Operation, error)
	DestroyApp(appID string) (Operation, error)
}

func NewAppDeployer(a AppDeployerConfig) (appDeployer AppDeployer, err error) {
	switch a.Type {
	case "marathon":
		return NewMarathonManager(a.Address)
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

func (m *marathonManager) NewApp() App {
	return &marathonApp{gomApp: marathon.NewDockerApplication()}
}

func (m *marathonManager) DeployApp(app App) (Operation, error) {
	marathonApp, ok := app.(*marathonApp)
	if !ok {
		return nil, errors.New("failed because didn't receive a marathon app")
	}
	gomApp, err := m.marathonClient.CreateApplication(marathonApp.gomApp)
	if err != nil {
		return nil, err
	}
	marathonDeploymentIDs := gomApp.DeploymentIDs()
	deploymentIDStrings := make([]string, len(marathonDeploymentIDs))
	for i, marathonDeploymentID := range marathonDeploymentIDs {
		deploymentIDStrings[i] = marathonDeploymentID.DeploymentID
	}
	op := &marathonDeployment{appID: gomApp.ID, deploymentIDs: deploymentIDStrings, marathonManager: *m}
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
