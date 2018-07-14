package nomad

import (
	"context"
	"fmt"

	"github.com/hashicorp/nomad/api"
	"github.com/pkg/errors"

	"git.corp.adobe.com/abramowi/hyperion/core"
	"git.corp.adobe.com/abramowi/hyperion/utils"
)

var ctx = context.TODO()

type manager struct {
	client     *api.Client
	jobsClient *api.Jobs
	url        string
}

// NewManager returns a Manager for Kubernetes.
func NewManager(url string) (*manager, error) {
	config := &api.Config{Address: url}
	client, err := api.NewClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "nomad.NewManager: api.NewClient failed")
	}
	return &manager{client: client, jobsClient: client.Jobs(), url: url}, nil
}

// Svcs returns info about all running services.
func (mgr *manager) Svcs() ([]core.Svc, error) {
	return nil, errors.New("nomad.manager.Svcs: Not implemented")
}

// SvcTasks returns info about the running tasks for a service.
func (mgr *manager) SvcTasks(svcCfg core.SvcCfg) ([]core.Task, error) {
	return nil, errors.New("nomad.manager.SvcTasks: Not implemented")
}

// Tasks returns info about all running tasks.
func (mgr *manager) Tasks() ([]core.Task, error) {
	return nil, errors.New("nomad.manager.Tasks: Not implemented")
}

// DeploySvc takes a SvcCfg and deploys it, returning an Operation.
func (mgr *manager) DeploySvc(svcCfg core.SvcCfg) (core.Operation, error) {
	job := getJob(svcCfg)
	jobRegisterResponse, writeMeta, err := mgr.jobsClient.Register(job, &api.WriteOptions{})
	fmt.Printf("*** jobRegisterResponse = %+v; writeMeta = %+v; err = %+v\n", jobRegisterResponse, writeMeta, err)
	if err != nil {
		return nil, errors.Wrap(err, "nomad.manager.DeploySvc: mgr.jobsClient.Register failed")
	}
	return nil, nil
}

// DestroySvc destroys a service.
func (mgr *manager) DestroySvc(svcID string) (core.Operation, error) {
	purge := true
	jobDeregisterResponse, writeMeta, err := mgr.jobsClient.Deregister(svcID, purge, &api.WriteOptions{})
	fmt.Printf("*** jobDeregisterResponse = %+v; writeMeta = %+v; err = %+v\n", jobDeregisterResponse, writeMeta, err)
	if err != nil {
		return nil, errors.Wrap(err, "nomad.manager.DestroySvc: mgr.jobsClient.Deregister failed")
	}
	return nil, err
}

func getJob(svcCfg core.SvcCfg) *api.Job {
	return &api.Job{
		ID:          utils.Sptr(svcCfg.ID),
		Name:        utils.Sptr(svcCfg.ID),
		Type:        utils.Sptr(api.JobTypeService),
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			&api.TaskGroup{
				Name:  utils.Sptr(svcCfg.ID),
				Count: &svcCfg.Count,
				Tasks: []*api.Task{
					&api.Task{
						Name:   svcCfg.ID,
						Driver: "docker",
						Config: map[string]interface{}{
							"image": svcCfg.Image,
						},
					},
				},
			},
		},
	}
}
