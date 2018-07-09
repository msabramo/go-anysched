package hyperion

import (
	"fmt"

	"git.corp.adobe.com/abramowi/hyperion/core"
	"git.corp.adobe.com/abramowi/hyperion/dockerswarm"
	"git.corp.adobe.com/abramowi/hyperion/kubernetes"
	"git.corp.adobe.com/abramowi/hyperion/marathon"
	"git.corp.adobe.com/abramowi/hyperion/nomad"
)

type AppDeployerConfig struct {
	Type    string // e.g.: "marathon", "kubernetes", etc.
	Address string // e.g.: "http://127.0.0.1:8080"
}

type AppDeployer interface {
	DeployApp(core.App) (core.Operation, error)
	DestroyApp(appID string) (core.Operation, error)
}

func NewAppDeployer(a AppDeployerConfig) (appDeployer AppDeployer, err error) {
	switch a.Type {
	case "marathon":
		return marathon.NewMarathonManager(a.Address)
	case "kubernetes":
		return kubernetes.NewK8sManager(a.Address)
	case "dockerswarm":
		return dockerswarm.NewDockerSwarmManager(a.Address)
	case "nomad":
		return nomad.NewNomadManager(a.Address)
	default:
		return nil, fmt.Errorf("Unknown type: %q", a.Type)
	}
}
