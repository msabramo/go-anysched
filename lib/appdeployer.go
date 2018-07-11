package hyperion

import (
	"fmt"

	"git.corp.adobe.com/abramowi/hyperion/lib/dockerswarm"
	"git.corp.adobe.com/abramowi/hyperion/lib/kubernetes"
	"git.corp.adobe.com/abramowi/hyperion/lib/marathon"
	"git.corp.adobe.com/abramowi/hyperion/lib/nomad"
)

// AppDeployerType represents the type of system we're managing apps on -- e.g.:
// Marathon, Kubernetes, etc.
type AppDeployerType string

const (
	AppDeployerTypeMarathon    AppDeployerType = "marathon"
	AppDeployerTypeKubernetes  AppDeployerType = "kubernetes"
	AppDeployerTypeDockerSwarm AppDeployerType = "dockerswarm"
	AppDeployerTypeNomad       AppDeployerType = "nomad"
)

// AppDeployerConfig contains config passed to the NewAppDeployer function
type AppDeployerConfig struct {
	Type    AppDeployerType // e.g.: "marathon", "kubernetes", etc.
	Address string          // e.g.: "http://127.0.0.1:8080"
}

type AppDeployer interface {
	DeployApp(App) (Operation, error)
	DestroyApp(appID string) (Operation, error)
}

func NewAppDeployer(a AppDeployerConfig) (appDeployer AppDeployer, err error) {
	switch a.Type {
	case AppDeployerTypeMarathon:
		return marathon.NewManager(a.Address)
	case AppDeployerTypeKubernetes:
		return kubernetes.NewManager(a.Address)
	case AppDeployerTypeDockerSwarm:
		return dockerswarm.NewManager(a.Address)
	case AppDeployerTypeNomad:
		return nomad.NewManager(a.Address)
	default:
		return nil, fmt.Errorf("Unknown type: %q", a.Type)
	}
}
