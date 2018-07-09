package nomad

import (
	"context"
	"fmt"

	"github.com/hashicorp/nomad/api"

	"git.corp.adobe.com/abramowi/hyperion/core"
)

var ctx = context.TODO()

type nomadManager struct {
	client *api.Client
	url    string
}

func NewNomadManager(url string) (*nomadManager, error) {
	config := &api.Config{Address: url}
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	manager := &nomadManager{
		client: client,
		url:    url,
	}
	return manager, nil
}

func Sptr(s string) *string {
	return &s
}

func (m *nomadManager) DeployApp(app core.App) (core.Operation, error) {
	job := &api.Job{
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
	q := &api.WriteOptions{}
	jobRegisterResponse, writeMeta, err := m.client.Jobs().Register(job, q)
	fmt.Printf("*** jobRegisterResponse = %+v; writeMeta = %+v; err = %+v\n", jobRegisterResponse, writeMeta, err)
	return nil, err
}

func (m *nomadManager) DestroyApp(appID string) (core.Operation, error) {
	purge := true
	q := &api.WriteOptions{}
	jobDeregisterResponse, writeMeta, err := m.client.Jobs().Deregister(appID, purge, q)
	fmt.Printf("*** jobDeregisterResponse = %+v; writeMeta = %+v; err = %+v\n", jobDeregisterResponse, writeMeta, err)
	return nil, err
}
