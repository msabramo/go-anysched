package dockerswarm

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/swarm"
	dockerclient "github.com/docker/docker/client"

	"github.com/msabramo/go-anysched"
)

var ctx = context.TODO()

type manager struct {
	client *dockerclient.Client
	url    string
}

func init() {
	anysched.RegisterManagerType("dockerswarm", NewManager)
}

// NewManager returns a Manager for Docker Swarm.
func NewManager(url string) (anysched.Manager, error) {
	client, err := dockerclient.NewEnvClient()
	if err != nil {
		return nil, errors.Wrap(err, "dockerswarm.NewManager: dockerclient.NewEnvClient failed")
	}
	return &manager{client: client, url: url}, nil
}

// Svcs returns info about all running services.
func (mgr *manager) Svcs() ([]anysched.Svc, error) {
	return nil, errors.New("dockerswarm.manager.Svcs: Not implemented")
}

// SvcTasks returns info about the running tasks for a service.
func (mgr *manager) SvcTasks(svcCfg anysched.SvcCfg) ([]anysched.Task, error) {
	return nil, errors.New("dockerswarm.manager.SvcTasks: Not implemented")
}

// Tasks returns info about all running tasks.
func (mgr *manager) Tasks() ([]anysched.Task, error) {
	return nil, errors.New("dockerswarm.manager.Tasks: Not implemented")
}

// DeploySvc takes a SvcCfg and deploys it, returning an Operation.
func (mgr *manager) DeploySvc(svcCfg anysched.SvcCfg) (anysched.Operation, error) {
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
func (mgr *manager) DestroySvc(svcID string) (anysched.Operation, error) {
	err := mgr.client.ServiceRemove(ctx, svcID)
	if err != nil {
		return nil, errors.Wrap(err, "dockerswarm.manager.DestroySvc: mgr.client.ServiceRemove failed")
	}
	return nil, nil
}
