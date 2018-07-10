package hyperion

import (
	"fmt"

	"git.corp.adobe.com/abramowi/hyperion/lib/core"
	"git.corp.adobe.com/abramowi/hyperion/lib/dockerswarm"
	"git.corp.adobe.com/abramowi/hyperion/lib/kubernetes"
	"git.corp.adobe.com/abramowi/hyperion/lib/marathon"
	"git.corp.adobe.com/abramowi/hyperion/lib/nomad"
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
		return marathon.NewManager(a.Address)
	case "kubernetes":
		return kubernetes.NewManager(a.Address)
	case "dockerswarm":
		return dockerswarm.NewManager(a.Address)
	case "nomad":
		return nomad.NewManager(a.Address)
	default:
		return nil, fmt.Errorf("Unknown type: %q", a.Type)
	}
}
