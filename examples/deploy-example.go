package main

import (
	"fmt"
	"os"

	"github.com/msabramo/go-anysched"
	_ "github.com/msabramo/go-anysched/managers/kubernetes"
)

func doStuff() error {
	managerConfig := anysched.ManagerConfig{Type: "kubernetes"}
	// or alternatively one of the following:
	//
	// managerConfig := anysched.ManagerConfig{
	// 	Type:    "marathon",
	// 	Address: "http://127.0.0.1:8080",
	// }
	// managerConfig := anysched.ManagerConfig{
	// 	Type:    "dockerswarm",
	// 	Address: "http://127.0.0.1:2377",
	// }
	// managerConfig := anysched.ManagerConfig{
	// 	Type:    "nomad",
	// 	Address: "http://127.0.0.1:4646",
	// }

	manager, err := anysched.NewManager(managerConfig)
	if err != nil {
		return err
	}
	svcCfg := anysched.SvcCfg{
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
		_, err = fmt.Fprintf(os.Stderr, "err: %s\n", err)
		if err != nil {
			panic(err)
		}
		os.Exit(1)
	}
}
