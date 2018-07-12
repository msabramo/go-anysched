package dockerswarm

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	dockerclient "github.com/docker/docker/client"
	"github.com/pkg/errors"

	"git.corp.adobe.com/abramowi/hyperion/core"
)

var ctx = context.TODO()

type manager struct {
	client *dockerclient.Client
	url    string
}

func NewManager(url string) (*manager, error) {
	client, err := dockerclient.NewEnvClient()
	if err != nil {
		return nil, errors.Wrap(err, "dockerswarm.NewManager: dockerclient.NewEnvClient failed")
	}
	return &manager{client: client, url: url}, nil
}

// AppTasks returns info about the running tasks for an app
func (mgr *manager) AppTasks(app core.App) (results []core.TaskInfo, err error) {
	return nil, errors.New("dockerswarm.manager.AppTasks: Not implemented")
}

// AllTasks returns info about all running tasks
func (mgr *manager) AllTasks() (results []core.TaskInfo, err error) {
	return nil, errors.New("dockerswarm.manager.AllTasks: Not implemented")
}

func (mgr *manager) DeployApp(app core.App) (core.Operation, error) {
	count := uint64(app.Count)
	service := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: app.ID,
		},
		Mode: swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: &count,
			},
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: swarm.ContainerSpec{
				Image: app.Image,
			},
		},
	}
	options := types.ServiceCreateOptions{}
	serviceCreateResponse, err := mgr.client.ServiceCreate(ctx, service, options)
	fmt.Printf("*** serviceCreateResponse = %+v; err = %+v\n", serviceCreateResponse, err)
	if err != nil {
		return nil, errors.Wrap(err, "dockerswarm.manager.DeployApp: mgr.client.ServiceCreate failed")
	}
	return nil, nil
}

func (mgr *manager) DestroyApp(appID string) (core.Operation, error) {
	err := mgr.client.ServiceRemove(ctx, appID)
	if err != nil {
		return nil, errors.Wrap(err, "dockerswarm.manager.DestroyApp: mgr.client.ServiceRemove failed")
	}
	return nil, nil
}
