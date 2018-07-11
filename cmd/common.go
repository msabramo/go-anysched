package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	hyperionlib "git.corp.adobe.com/abramowi/hyperion/lib"
)

func Manager() hyperionlib.Manager {
	managerConfig := hyperionlib.ManagerConfig{
		Type:    hyperionlib.ManagerTypeKubernetes,
		Address: "kubeconfig",
	}
	// or alternatively one of the following:
	//
	// managerConfig := ManagerConfig{
	// 	Type:    hyperionlib.ManagerTypeMarathon,
	// 	Address: "http://127.0.0.1:8080",
	// }
	// managerConfig := ManagerConfig{
	// 	Type:    hyperionlib.ManagerTypeDockerSwarm,
	// 	Address: "http://127.0.0.1:2377",
	// }
	// managerConfig := ManagerConfig{
	// 	Type:    hyperionlib.ManagerTypeNomad
	// 	Address: "http://127.0.0.1:4646",
	// }

	manager, err := hyperionlib.NewManager(managerConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	return manager
}

func WaitForCompletion(ctx context.Context, operation hyperionlib.Operation) error {
	if asyncOperation, ok := operation.(hyperionlib.AsyncOperation); ok && asyncOperation != nil {
		return asyncOperation.Wait(ctx, 15*time.Second)
	}
	return nil
}
