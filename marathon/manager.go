package marathon

import (
	"time"

	goMarathon "github.com/gambol99/go-marathon"
	"github.com/pkg/errors"

	"git.corp.adobe.com/abramowi/hyperion/core"
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

// AppTasks returns info about the running tasks for an app
func (mgr *manager) AppTasks(app core.App) (results []core.TaskInfo, err error) {
	// return nil, errors.New("marathon.manager.AppTasks: Not implemented")
	tasks, err := mgr.goMarathonClient.Tasks(app.ID)
	if err != nil {
		return nil, errors.Wrap(err, "marathon.manager.AppTasks: goMarathonClient.Tasks failed")
	}

	tasksSlice := tasks.Tasks
	results = make([]core.TaskInfo, len(tasksSlice))
	for i, task := range tasksSlice {
		// fmt.Printf("*** task = %+v\n", task)
		taskStartTime := new(time.Time)
		if task.StartedAt != "" {
			startTime, err := time.Parse(time.RFC3339, task.StartedAt)
			if err != nil {
				return nil, errors.Wrapf(err, "marathon.manager.AllTasks: time.Parse failed for %s", task.StartedAt)
			}
			*taskStartTime = startTime
		}
		ipAddresses := make([]string, len(task.IPAddresses))
		for i, addr := range task.IPAddresses {
			ipAddresses[i] = addr.IPAddress
		}
		results[i] = core.TaskInfo{
			Name:         task.ID,
			AppID:        task.AppID,
			HostName:     task.Host,
			IPAddresses:  ipAddresses,
			Ports:        task.Ports,
			ServicePorts: task.ServicePorts,
			MesosSlaveID: task.SlaveID,
			StartTime:    taskStartTime,
			State:        task.State,
			Version:      task.Version,
		}
	}
	return results, nil
}

// AllTasks returns info about all running tasks
func (mgr *manager) AllTasks() (results []core.TaskInfo, err error) {
	var opts *goMarathon.AllTasksOpts
	tasks, err := mgr.goMarathonClient.AllTasks(opts)
	if err != nil {
		return nil, errors.Wrap(err, "marathon.manager.AllTasks: goMarathonClient.AllTasks failed")
	}

	tasksSlice := tasks.Tasks
	results = make([]core.TaskInfo, len(tasksSlice))
	for i, task := range tasksSlice {
		// fmt.Printf("*** task = %+v\n", task)
		taskStartTime := new(time.Time)
		if task.StartedAt != "" {
			startTime, err := time.Parse(time.RFC3339, task.StartedAt)
			if err != nil {
				return nil, errors.Wrapf(err, "marathon.manager.AllTasks: time.Parse failed for %s", task.StartedAt)
			}
			*taskStartTime = startTime
		}
		ipAddresses := make([]string, len(task.IPAddresses))
		for i, addr := range task.IPAddresses {
			ipAddresses[i] = addr.IPAddress
		}
		results[i] = core.TaskInfo{
			Name:         task.ID,
			AppID:        task.AppID,
			HostName:     task.Host,
			IPAddresses:  ipAddresses,
			Ports:        task.Ports,
			ServicePorts: task.ServicePorts,
			MesosSlaveID: task.SlaveID,
			StartTime:    taskStartTime,
			State:        task.State,
			Version:      task.Version,
		}
	}
	return results, nil
}

func (mgr *manager) DeployApp(app core.App) (core.Operation, error) {
	gomApp, err := mgr.goMarathonClient.CreateApplication(goMarathonApp(app))
	if err != nil {
		return nil, errors.Wrap(err, "marathon.manager.DeployApp: goMarathonClient.CreateApplication failed")
	}
	return mgr.newDeploymentFromGoMarathonApp(gomApp), nil
}

func (mgr *manager) DestroyApp(appID string) (core.Operation, error) {
	force := false
	marathonDeploymentID, err := mgr.goMarathonClient.DeleteApplication(appID, force)
	if err != nil {
		return nil, err
	}
	op := &marathonDeploymentOperation{
		appID:           appID,
		deploymentIDs:   []string{marathonDeploymentID.DeploymentID},
		manager:         mgr,
		timeoutDuration: 60 * time.Second,
	}
	return op, err
}

func (mgr *manager) newDeploymentFromGoMarathonApp(gomApp *goMarathon.Application) *marathonDeploymentOperation {
	return &marathonDeploymentOperation{
		appID:           gomApp.ID,
		deploymentIDs:   deploymentIDs(gomApp),
		manager:         mgr,
		timeoutDuration: 60 * time.Second,
	}
}

func goMarathonApp(app core.App) (gomApp *goMarathon.Application) {
	gomApp = goMarathon.NewDockerApplication()
	gomApp.ID = app.ID
	gomApp.Container.Docker.Bridged()
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
