package marathon

import (
	goMarathon "github.com/gambol99/go-marathon"

	"git.corp.adobe.com/abramowi/hyperion/lib/core"
)

type marathonManager struct {
	goMarathonClient goMarathon.Marathon
	url              string
}

func NewMarathonManager(url string) (*marathonManager, error) {
	config := goMarathon.NewDefaultConfig()
	config.URL = url
	client, err := goMarathon.NewClient(config)
	if err != nil {
		return nil, err
	}
	m := &marathonManager{goMarathonClient: client, url: url}
	return m, nil
}

func (m *marathonManager) goMarathonApp(app core.App) (gomApp *goMarathon.Application) {
	gomApp = goMarathon.NewDockerApplication()
	gomApp.ID = app.ID
	gomApp.Container.Docker.Container(app.Image)
	gomApp.Count(app.Count)
	return gomApp
}

func (m *marathonManager) deploymentIDs(gomApp *goMarathon.Application) (deploymentIDs []string) {
	marathonDeploymentIDs := gomApp.DeploymentIDs()
	deploymentIDs = make([]string, len(marathonDeploymentIDs))
	for i, marathonDeploymentID := range marathonDeploymentIDs {
		deploymentIDs[i] = marathonDeploymentID.DeploymentID
	}
	return deploymentIDs
}

func (m *marathonManager) DeployApp(app core.App) (core.Operation, error) {
	gomApp, err := m.goMarathonClient.CreateApplication(m.goMarathonApp(app))
	if err != nil {
		return nil, err
	}
	op := &marathonDeployment{
		appID:           gomApp.ID,
		deploymentIDs:   m.deploymentIDs(gomApp),
		marathonManager: *m,
	}
	return op, err
}

func (m *marathonManager) DestroyApp(appID string) (core.Operation, error) {
	force := false
	marathonDeploymentID, err := m.goMarathonClient.DeleteApplication(appID, force)
	if err != nil {
		return nil, err
	}
	op := &marathonDeployment{appID: appID, deploymentIDs: []string{marathonDeploymentID.DeploymentID}, marathonManager: *m}
	return op, err
}
