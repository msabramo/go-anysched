package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"git.corp.adobe.com/abramowi/hyperion"
)

func Manager() hyperion.Manager {
	managerConfig := hyperion.ManagerConfig{Type: hyperion.ManagerTypeKubernetes}
	// or alternatively one of the following:
	//
	// managerConfig := hyperonlib.ManagerConfig{
	// 	Type:    hyperion.ManagerTypeMarathon,
	// 	Address: "http://127.0.0.1:8080",
	// }
	// managerConfig := hyperonlib.ManagerConfig{
	// 	Type:    hyperion.ManagerTypeDockerSwarm,
	// 	Address: "http://127.0.0.1:2377",
	// }
	// managerConfig := hyperonlib.ManagerConfig{
	// 	Type:    hyperion.ManagerTypeNomad
	// 	Address: "http://127.0.0.1:4646",
	// }

	manager, err := hyperion.NewManager(managerConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	return manager
}

func WaitForCompletion(ctx context.Context, operation hyperion.Operation) error {
	if asyncOperation, ok := operation.(hyperion.AsyncOperation); ok && asyncOperation != nil {
		return asyncOperation.Wait(ctx, 15*time.Second)
	}
	return nil
}
