package hyperion

import (
	"context"
	"errors"
	"fmt"

	marathon "github.com/gambol99/go-marathon"
)

type App interface {
	SetCount(int) App
	SetDockerImage(string) App
	SetID(string) App
}

type Operation interface {
	// Async() AsyncOperation // returns an AsyncOperation if async; else nil
}

type AsyncOperation interface {
	Wait(context.Context) error // Wait for async operation to finish and return error or nil
}

type marathonDeployment struct {
	deploymentID string
	appID        string
}

func (d *marathonDeployment) Wait(ctx context.Context) error {
	fmt.Printf("Wait() called with d = %+v\n", d)
	return nil
}

/*
func (op *marathonOperation) Async() AsyncOperation {
	if asyncOp, ok := op.Operation.(AsyncOperation); ok && asyncOp != nil {
		return asyncOp
	}
	return nil
}
*/

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
		return NewMarathonManager()
	default:
		return nil, fmt.Errorf("Unknown type: %q", a.Type)
	}
}

type marathonManager struct {
	marathonClient marathon.Marathon
}

func NewMarathonManager() (*marathonManager, error) {
	config := marathon.NewDefaultConfig()
	config.URL = "http://127.0.0.1:8080"
	client, err := marathon.NewClient(config)
	if err != nil {
		return nil, err
	}
	m := &marathonManager{marathonClient: client}
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
	op := &marathonDeployment{appID: gomApp.ID}
	return op, err
}

func (m *marathonManager) DestroyApp(appID string) (Operation, error) {
	force := false
	marathonDeploymentID, err := m.marathonClient.DeleteApplication(appID, force)
	fmt.Printf("marathonDeploymentID = %+v\n", marathonDeploymentID)
	if err != nil {
		return nil, err
	}
	op := &marathonDeployment{appID: appID}
	return op, err
}
