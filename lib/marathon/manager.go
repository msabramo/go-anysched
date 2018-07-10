package marathon

import (
	goMarathon "github.com/gambol99/go-marathon"
	"github.com/pkg/errors"

	"git.corp.adobe.com/abramowi/hyperion/lib/core"
)

type manager struct {
	goMarathonClient goMarathon.Marathon
	url              string
}

func NewManager(url string) (*manager, error) {
	config := goMarathon.NewDefaultConfig()
	config.URL = url
	client, err := goMarathon.NewClient(config)
	if err != nil {
		return nil, err
	}
	mgr := &manager{goMarathonClient: client, url: url}
	return mgr, nil
}

func (mgr *manager) DeployApp(app core.App) (core.Operation, error) {
	gomApp, err := mgr.goMarathonClient.CreateApplication(goMarathonApp(app))
	if err != nil {
		return nil, errors.Wrap(err, "marathon.manager.DeployApp: goMarathonClient.CreateApplication failed")
	}
	return mgr.newDeploymentFromGoMarathonApp(gomApp), nil
}

func (m *manager) DestroyApp(appID string) (core.Operation, error) {
	force := false
	marathonDeploymentID, err := m.goMarathonClient.DeleteApplication(appID, force)
	if err != nil {
		return nil, err
	}
	op := &marathonDeployment{appID: appID, deploymentIDs: []string{marathonDeploymentID.DeploymentID}, manager: *m}
	return op, err
}

func (mgr *manager) newDeploymentFromGoMarathonApp(gomApp *goMarathon.Application) *marathonDeployment {
	return &marathonDeployment{
		appID:         gomApp.ID,
		deploymentIDs: deploymentIDs(gomApp),
		manager:       *mgr,
	}
}

func goMarathonApp(app core.App) (gomApp *goMarathon.Application) {
	gomApp = goMarathon.NewDockerApplication()
	gomApp.ID = app.ID
	gomApp.Container.Docker.Container(app.Image)
	gomApp.Count(app.Count)
	return gomApp
}

func deploymentIDs(gomApp *goMarathon.Application) (deploymentIDs []string) {
	marathonDeploymentIDs := gomApp.DeploymentIDs()
	deploymentIDs = make([]string, len(marathonDeploymentIDs))
	for i, marathonDeploymentID := range marathonDeploymentIDs {
		deploymentIDs[i] = marathonDeploymentID.DeploymentID
	}
	return deploymentIDs
}
