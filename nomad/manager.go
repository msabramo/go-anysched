package nomad

import (
	"context"
	"fmt"

	"github.com/hashicorp/nomad/api"
	"github.com/pkg/errors"

	"git.corp.adobe.com/abramowi/hyperion/core"
)

var ctx = context.TODO()

type manager struct {
	client     *api.Client
	jobsClient *api.Jobs
	url        string
}

func NewManager(url string) (*manager, error) {
	config := &api.Config{Address: url}
	client, err := api.NewClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "nomad.NewManager: api.NewClient failed")
	}
	return &manager{client: client, jobsClient: client.Jobs(), url: url}, nil
}

func Sptr(s string) *string {
	return &s
}

// AllApps returns info about all running apps
func (mgr *manager) AllApps() (results []core.AppInfo, err error) {
	return nil, errors.New("nomad.manager.AllApps: Not implemented")
}

// AppTasks returns info about the running tasks for an app
func (mgr *manager) AppTasks(app core.App) (results []core.TaskInfo, err error) {
	return nil, errors.New("nomad.manager.AppTasks: Not implemented")
}

// AllTasks returns info about all running tasks
func (mgr *manager) AllTasks() (results []core.TaskInfo, err error) {
	return nil, errors.New("nomad.manager.AllTasks: Not implemented")
}

func (mgr *manager) DeployApp(app core.App) (core.Operation, error) {
	job := getJob(app)
	jobRegisterResponse, writeMeta, err := mgr.jobsClient.Register(job, &api.WriteOptions{})
	fmt.Printf("*** jobRegisterResponse = %+v; writeMeta = %+v; err = %+v\n", jobRegisterResponse, writeMeta, err)
	if err != nil {
		return nil, errors.Wrap(err, "nomad.manager.DeployApp: mgr.jobsClient.Register failed")
	}
	return nil, nil
}

func (mgr *manager) DestroyApp(appID string) (core.Operation, error) {
	purge := true
	jobDeregisterResponse, writeMeta, err := mgr.jobsClient.Deregister(appID, purge, &api.WriteOptions{})
	fmt.Printf("*** jobDeregisterResponse = %+v; writeMeta = %+v; err = %+v\n", jobDeregisterResponse, writeMeta, err)
	if err != nil {
		return nil, errors.Wrap(err, "nomad.manager.DestroyApp: mgr.jobsClient.Deregister failed")
	}
	return nil, err
}

func getJob(app core.App) *api.Job {
	return &api.Job{
		ID:          Sptr(app.ID),
		Name:        Sptr(app.ID),
		Type:        Sptr(api.JobTypeService),
		Datacenters: []string{"dc1"},
		TaskGroups: []*api.TaskGroup{
			&api.TaskGroup{
				Name:  Sptr(app.ID),
				Count: &app.Count,
				Tasks: []*api.Task{
					&api.Task{
						Name:   app.ID,
						Driver: "docker",
						Config: map[string]interface{}{
							"image": app.Image,
						},
					},
				},
			},
		},
	}
}
