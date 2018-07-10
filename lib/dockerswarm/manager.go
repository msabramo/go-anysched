package dockerswarm

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	dockerclient "github.com/docker/docker/client"

	"git.corp.adobe.com/abramowi/hyperion/lib/core"
)

var ctx = context.TODO()

type dockerSwarmManager struct {
	client *dockerclient.Client
	url    string
}

func NewDockerSwarmManager(url string) (*dockerSwarmManager, error) {
	client, err := dockerclient.NewEnvClient()
	if err != nil {
		return nil, err
	}
	manager := &dockerSwarmManager{
		client: client,
		url:    url,
	}
	return manager, nil
}

func (m *dockerSwarmManager) DeployApp(app core.App) (core.Operation, error) {
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
	serviceCreateResponse, err := m.client.ServiceCreate(ctx, service, options)
	fmt.Printf("*** serviceCreateResponse = %+v; err = %+v\n", serviceCreateResponse, err)
	return nil, err
}

func (m *dockerSwarmManager) DestroyApp(appID string) (core.Operation, error) {
	err := m.client.ServiceRemove(ctx, appID)
	return nil, err
}
