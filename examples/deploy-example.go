package main

import (
	"fmt"
	"os"

	"git.corp.adobe.com/abramowi/hyperion"
	_ "git.corp.adobe.com/abramowi/hyperion/kubernetes"
)

func doStuff() error {
	managerConfig := hyperion.ManagerConfig{Type: "kubernetes"}
	// or alternatively one of the following:
	//
	// managerConfig := hyperion.ManagerConfig{
	// 	Type:    "marathon",
	// 	Address: "http://127.0.0.1:8080",
	// }
	// managerConfig := hyperion.ManagerConfig{
	// 	Type:    "dockerswarm",
	// 	Address: "http://127.0.0.1:2377",
	// }
	// managerConfig := hyperion.ManagerConfig{
	// 	Type:    "nomad",
	// 	Address: "http://127.0.0.1:4646",
	// }

	manager, err := hyperion.NewManager(managerConfig)
	if err != nil {
		return err
	}
	svcCfg := hyperion.SvcCfg{
		ID:    "my-svc-id",
		Image: "citizenstig/httpbin",
		Count: 4,
	}
	operation, err := manager.DeploySvc(svcCfg)
	if err != nil {
		return err
	}
	fmt.Printf("operation = %+v\n", operation)
	return nil
}

func main() {
	err := doStuff()
	if err != nil {
		fmt.Fprintf(os.Stderr, "err: %s\n", err)
		os.Exit(1)
	}
}
