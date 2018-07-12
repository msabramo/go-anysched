package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"git.corp.adobe.com/abramowi/hyperion"
	yaml "gopkg.in/yaml.v2"
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

func output(w io.Writer, data interface{}, format string, outputTable func(w io.Writer, data interface{}) error) error {
	switch format {
	case "yaml":
		return outputYAML(w, data)
	case "json":
		return outputJSON(w, data)
	case "table":
		return outputTable(w, data)
	default:
		return fmt.Errorf("unknown output format type: %q", format)
	}
}

func outputYAML(w io.Writer, data interface{}) error {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return err
	}
	w.Write(bytes)
	return nil
}

func outputJSON(w io.Writer, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	w.Write(bytes)
	return nil
}
