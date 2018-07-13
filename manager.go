package hyperion

import (
	"fmt"

	"git.corp.adobe.com/abramowi/hyperion/core"
	"git.corp.adobe.com/abramowi/hyperion/dockerswarm"
	"git.corp.adobe.com/abramowi/hyperion/kubernetes"
	"git.corp.adobe.com/abramowi/hyperion/marathon"
	"git.corp.adobe.com/abramowi/hyperion/nomad"
)

type App = core.App
type AppInfo = core.AppInfo
type Operation = core.Operation
type TaskInfo = core.TaskInfo

// ManagerType represents the type of system we're managing apps on -- e.g.:
// Marathon, Kubernetes, etc.
type ManagerType string

const (
	// ManagerTypeMarathon is a const ManagerType (string) for Marathon.
	ManagerTypeMarathon ManagerType = "marathon"

	// ManagerTypeKubernetes is a const ManagerType (string) for Kubernetes.
	ManagerTypeKubernetes ManagerType = "kubernetes"

	// ManagerTypeDockerSwarm is a const ManagerType (string) for Docker Swarm.
	ManagerTypeDockerSwarm ManagerType = "dockerswarm"

	// ManagerTypeNomad is a const ManagerType (string) for Nomad.
	ManagerTypeNomad ManagerType = "nomad"
)

// ManagerTypes is a slice with valid manager types
var ManagerTypes = [...]ManagerType{
	ManagerTypeMarathon,
	ManagerTypeKubernetes,
	ManagerTypeDockerSwarm,
	ManagerTypeNomad,
}

// ManagerConfig is a struct containing configuration info that a user passes to
// the NewManager function.
type ManagerConfig struct {
	Type    ManagerType // e.g.: "marathon", "kubernetes", etc.
	Address string      // e.g.: "http://127.0.0.1:8080"
}

type AllAppsGetter interface {
	// AllApps returns info about the running apps.
	AllApps() (results []AppInfo, err error)
}

type AppTasksGetter interface {
	// AppTasks returns info about the running tasks for an app.
	AppTasks(app core.App) (results []TaskInfo, err error)
}

type AllTasksGetter interface {
	// AllTasks returns info about all running tasks.
	AllTasks() (results []TaskInfo, err error)
}

type Deployer interface {
	DeployApp(App) (Operation, error)
}

type Destroyer interface {
	DestroyApp(appID string) (Operation, error)
}

// Manager is an interface that is composed of various other more fine-grained
// interfaces.
type Manager interface {
	AllAppsGetter
	AllTasksGetter
	AppTasksGetter
	Deployer
	Destroyer
}

// NewManager takes a ManagerConfig and returns a specific type of Manager for
// the scheduler that the user requested (e.g.: Kubernetes, Marathon, etc.).
func NewManager(managerConfig ManagerConfig) (manager Manager, err error) {
	switch managerConfig.Type {
	case ManagerTypeMarathon:
		return marathon.NewManager(managerConfig.Address)
	case ManagerTypeKubernetes:
		return kubernetes.NewManager(managerConfig.Address)
	case ManagerTypeDockerSwarm:
		return dockerswarm.NewManager(managerConfig.Address)
	case ManagerTypeNomad:
		return nomad.NewManager(managerConfig.Address)
	default:
		return nil, unknownAppManagerTypeError(managerConfig.Type)
	}
}

func unknownAppManagerTypeError(appManagerType ManagerType) error {
	return fmt.Errorf("Unknown app manager type: %q. Valid options are: %+v", appManagerType, ManagerTypes)
}
