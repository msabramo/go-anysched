package marathon

import (
	"net/url"
	"time"

	goMarathon "github.com/gambol99/go-marathon"
	"github.com/pkg/errors"

	"git.corp.adobe.com/abramowi/hyperion"
	"git.corp.adobe.com/abramowi/hyperion/core"
)

var (
	goMarathonDefaultAllTasksOpts *goMarathon.AllTasksOpts // = nil
)

var (
	goMarathonEmbedTasks = url.Values{"embed": []string{"apps.tasks"}}
)

type manager struct {
	goMarathonClient goMarathon.Marathon
	url              string
}

func init() {
	hyperion.RegisterManagerType("marathon", NewManager)
}

// NewManager returns a Manager for Marathon.
func NewManager(url string) (hyperion.Manager, error) {
	config := goMarathon.NewDefaultConfig()
	config.URL = url
	client, err := goMarathon.NewClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "marathon.NewManager: goMarathon.NewClient failed")
	}
	mgr := &manager{goMarathonClient: client, url: url}
	return mgr, nil
}

// Svcs returns info about all running services.
func (mgr *manager) Svcs() ([]core.Svc, error) {
	goMarathonAppsStruct, err := mgr.goMarathonClient.Applications(goMarathonEmbedTasks)
	if err != nil {
		return nil, errors.Wrap(err, "marathon.manager.Svcs: goMarathonClient.Svcs failed")
	}
	goMarathonAppsSlice := goMarathonAppsStruct.Apps
	svcs := make([]core.Svc, len(goMarathonAppsSlice))
	for i, goMarathonApp := range goMarathonAppsSlice {
		svcs[i] = svcFromMarathonApp(goMarathonApp)
	}
	return svcs, nil
}

func svcFromMarathonApp(goMarathonApp goMarathon.Application) core.Svc {
	return core.Svc{
		ID:             goMarathonApp.ID,
		TasksRunning:   &goMarathonApp.TasksRunning,
		TasksHealthy:   &goMarathonApp.TasksHealthy,
		TasksUnhealthy: &goMarathonApp.TasksUnhealthy,
	}
}

// SvcTasks returns info about the running tasks for a service.
func (mgr *manager) SvcTasks(svcCfg core.SvcCfg) ([]core.Task, error) {
	goMarathonTasksStruct, err := mgr.goMarathonClient.Tasks(svcCfg.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "marathon.manager.SvcTasks: goMarathonClient.Tasks(%q) failed", svcCfg.ID)
	}

	goMarathonTasksSlice := goMarathonTasksStruct.Tasks
	ourTasks := make([]core.Task, len(goMarathonTasksSlice))
	for i, goMarathonTask := range goMarathonTasksSlice {
		ourTask, err := ourTaskForGoMarathonTask(goMarathonTask)
		if err != nil {
			return nil, err
		}
		ourTasks[i] = *ourTask
	}
	return ourTasks, nil
}

// Tasks returns info about all running tasks.
func (mgr *manager) Tasks() ([]core.Task, error) {
	goMarathonTasksStruct, err := mgr.goMarathonClient.AllTasks(goMarathonDefaultAllTasksOpts)
	if err != nil {
		return nil, errors.Wrap(err, "marathon.manager.Tasks: goMarathonClient.AllTasks failed")
	}

	goMarathonTasksSlice := goMarathonTasksStruct.Tasks
	ourTasks := make([]core.Task, len(goMarathonTasksSlice))
	for i, goMarathonTask := range goMarathonTasksSlice {
		ourTask, err := ourTaskForGoMarathonTask(goMarathonTask)
		if err != nil {
			return nil, err
		}
		ourTasks[i] = *ourTask
	}
	return ourTasks, nil
}

func ourTaskForGoMarathonTask(goMarathonTask goMarathon.Task) (*core.Task, error) {
	taskStageTime, err := parseMarathonTime(goMarathonTask.StagedAt)
	if err != nil {
		return nil, err
	}
	taskStartTime, err := parseMarathonTime(goMarathonTask.StartedAt)
	if err != nil {
		return nil, err
	}
	ipAddresses := make([]string, len(goMarathonTask.IPAddresses))
	for i, goMarathonIPAddressStruct := range goMarathonTask.IPAddresses {
		ipAddresses[i] = goMarathonIPAddressStruct.IPAddress
	}
	return &core.Task{
		Name:         goMarathonTask.ID,
		AppID:        goMarathonTask.AppID,
		HostName:     goMarathonTask.Host,
		IPAddresses:  ipAddresses,
		Ports:        goMarathonTask.Ports,
		ServicePorts: goMarathonTask.ServicePorts,
		MesosSlaveID: goMarathonTask.SlaveID,
		StageTime:    taskStageTime,
		StartTime:    taskStartTime,
		State:        goMarathonTask.State,
		Version:      goMarathonTask.Version,
	}, nil
}

// parseMarathonTime takes a timestamp string from Marathon (formatted as
// RFC3339) and parses it into a *time.Time.
func parseMarathonTime(marathonTime string) (*time.Time, error) {
	if marathonTime != "" {
		t, err := time.Parse(time.RFC3339, marathonTime)
		if err != nil {
			return nil, errors.Wrapf(err,
				"marathon.manager.parseMarathonTime: time.Parse failed for: %s", marathonTime)
		}
		return &t, nil
	}
	return nil, nil
}

// DeploySvc takes a SvcCfg and deploys it, returning an Operation.
func (mgr *manager) DeploySvc(svcCfg core.SvcCfg) (core.Operation, error) {
	goMarathonApp, err := mgr.goMarathonClient.CreateApplication(goMarathonApp(svcCfg))
	if err != nil {
		return nil, errors.Wrap(err, "marathon.manager.DeploySvc: goMarathonClient.CreateApplication failed")
	}
	return mgr.newDeploymentFromGoMarathonApp(goMarathonApp), nil
}

// DestroySvc destroys a service.
func (mgr *manager) DestroySvc(svcID string) (core.Operation, error) {
	force := false
	marathonDeploymentID, err := mgr.goMarathonClient.DeleteApplication(svcID, force)
	if err != nil {
		return nil, errors.Wrap(err, "marathon.manager.DestroySvc: goMarathonClient.DeleteApplication failed")
	}
	op := &deployment{
		svcID: svcID,
		marathonDeploymentIDs: []string{marathonDeploymentID.DeploymentID},
		manager:               mgr,
		timeoutDuration:       60 * time.Second,
	}
	return op, err
}

func (mgr *manager) newDeploymentFromGoMarathonApp(goMarathonApp *goMarathon.Application) *deployment {
	return &deployment{
		svcID: goMarathonApp.ID,
		marathonDeploymentIDs: marathonDeploymentIDs(goMarathonApp),
		manager:               mgr,
		timeoutDuration:       60 * time.Second,
	}
}

func goMarathonApp(svcCfg core.SvcCfg) *goMarathon.Application {
	goMarathonApp := goMarathon.NewDockerApplication()
	goMarathonApp.ID = svcCfg.ID
	goMarathonApp.Container.Docker.Bridged()
	goMarathonApp.Container.Docker.Container(svcCfg.Image)
	goMarathonApp.Count(svcCfg.Count)
	return goMarathonApp
}

func marathonDeploymentIDs(goMarathonApp *goMarathon.Application) (marathonDeploymentIDs []string) {
	marathonDeploymentIDStructs := goMarathonApp.DeploymentIDs()
	marathonDeploymentIDs = make([]string, len(marathonDeploymentIDStructs))
	for i, marathonDeploymentIDStruct := range marathonDeploymentIDStructs {
		marathonDeploymentIDs[i] = marathonDeploymentIDStruct.DeploymentID
	}
	return marathonDeploymentIDs
}
