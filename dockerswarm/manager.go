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

// NewManager returns a Manager for Docker Swarm.
func NewManager(url string) (*manager, error) {
	client, err := dockerclient.NewEnvClient()
	if err != nil {
		return nil, errors.Wrap(err, "dockerswarm.NewManager: dockerclient.NewEnvClient failed")
	}
	return &manager{client: client, url: url}, nil
}

// Svcs returns info about all running services.
func (mgr *manager) Svcs() ([]core.Svc, error) {
	return nil, errors.New("dockerswarm.manager.Svcs: Not implemented")
}

// SvcTasks returns info about the running tasks for a service.
func (mgr *manager) SvcTasks(svcCfg core.SvcCfg) ([]core.Task, error) {
	return nil, errors.New("dockerswarm.manager.SvcTasks: Not implemented")
}

// Tasks returns info about all running tasks.
func (mgr *manager) Tasks() ([]core.Task, error) {
	return nil, errors.New("dockerswarm.manager.Tasks: Not implemented")
}

// DeploySvc takes a SvcCfg and deploys it, returning an Operation.
func (mgr *manager) DeploySvc(svcCfg core.SvcCfg) (core.Operation, error) {
	count := uint64(svcCfg.Count)
	service := swarm.ServiceSpec{
		Annotations: swarm.Annotations{
			Name: svcCfg.ID,
		},
		Mode: swarm.ServiceMode{
			Replicated: &swarm.ReplicatedService{
				Replicas: &count,
			},
		},
		TaskTemplate: swarm.TaskSpec{
			ContainerSpec: swarm.ContainerSpec{
				Image: svcCfg.Image,
			},
		},
	}
	options := types.ServiceCreateOptions{}
	serviceCreateResponse, err := mgr.client.ServiceCreate(ctx, service, options)
	fmt.Printf("*** serviceCreateResponse = %+v; err = %+v\n", serviceCreateResponse, err)
	if err != nil {
		return nil, errors.Wrap(err, "dockerswarm.manager.DeploySvc: mgr.client.ServiceCreate failed")
	}
	return nil, nil
}

// DestroySvc destroys a service.
func (mgr *manager) DestroySvc(svcID string) (core.Operation, error) {
	err := mgr.client.ServiceRemove(ctx, svcID)
	if err != nil {
		return nil, errors.Wrap(err, "dockerswarm.manager.DestroySvc: mgr.client.ServiceRemove failed")
	}
	return nil, nil
}
