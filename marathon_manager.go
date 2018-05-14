package hyperion

import (
	"fmt"

	marathon "github.com/gambol99/go-marathon"
)

type App interface{}

type Operation interface {
	Error() error
}

type DeployAppRequest struct {
	app App
}

type DestroyAppRequest struct {
	app App
}

type AppDeployer interface {
	NewApp() App
	DeployApp(DeployAppRequest) Operation
	DestroyApp(DestroyAppRequest) Operation
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

func (m *marathonManager) NewApp() *marathonApp {
	return &marathonApp{gomApp: marathon.NewDockerApplication()}
}

func (m *marathonManager) CreateApplication(app *marathonApp) (*marathonApp, error) {
	gomApp, err := m.marathonClient.CreateApplication(app.gomApp)
	if err != nil {
		return nil, err
	}
	return &marathonApp{gomApp: gomApp}, nil
}

func (m *marathonManager) DeleteApplication(deleteRequest *marathonAppDeleteRequest) error {
	force := false
	marathonDeploymentID, err := m.marathonClient.DeleteApplication(deleteRequest.appID, force)
	fmt.Printf("marathonDeploymentID = %+v\n", marathonDeploymentID)
	if err != nil {
		return err
	}
	return nil
}
