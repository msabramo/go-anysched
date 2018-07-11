package hyperion

import (
	"fmt"

	"git.corp.adobe.com/abramowi/hyperion/lib/core"
	"git.corp.adobe.com/abramowi/hyperion/lib/dockerswarm"
	"git.corp.adobe.com/abramowi/hyperion/lib/kubernetes"
	"git.corp.adobe.com/abramowi/hyperion/lib/marathon"
	"git.corp.adobe.com/abramowi/hyperion/lib/nomad"
)

type App = core.App
type Operation = core.Operation
type AsyncOperation = core.AsyncOperation

// ManagerType represents the type of system we're managing apps on -- e.g.:
// Marathon, Kubernetes, etc.
type ManagerType string

const (
	ManagerTypeMarathon    ManagerType = "marathon"
	ManagerTypeKubernetes  ManagerType = "kubernetes"
	ManagerTypeDockerSwarm ManagerType = "dockerswarm"
	ManagerTypeNomad       ManagerType = "nomad"
)

var ManagerTypes = [...]ManagerType{
	ManagerTypeMarathon,
	ManagerTypeKubernetes,
	ManagerTypeDockerSwarm,
	ManagerTypeNomad,
}

// ManagerConfig contains config passed to the NewManager function
type ManagerConfig struct {
	Type    ManagerType // e.g.: "marathon", "kubernetes", etc.
	Address string      // e.g.: "http://127.0.0.1:8080"
}

type Deployer interface {
	DeployApp(App) (Operation, error)
}

type Destroyer interface {
	DestroyApp(appID string) (Operation, error)
}

type Manager interface {
	Deployer
	Destroyer
}

func NewManager(a ManagerConfig) (manager Manager, err error) {
	switch a.Type {
	case ManagerTypeMarathon:
		return marathon.NewManager(a.Address)
	case ManagerTypeKubernetes:
		return kubernetes.NewManager(a.Address)
	case ManagerTypeDockerSwarm:
		return dockerswarm.NewManager(a.Address)
	case ManagerTypeNomad:
		return nomad.NewManager(a.Address)
	default:
		return nil, fmt.Errorf("Unknown app manager type: %q. Valid options are: %+v", a.Type, ManagerTypes)
	}
}
