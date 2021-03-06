package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/lithammer/dedent"
	yaml "gopkg.in/yaml.v2"
)

func die(format string, a ...interface{}) {
	if _, err := fmt.Fprintf(os.Stderr, dedent.Dedent(format)+"\n", a...); err != nil {
		panic(err)
	}
	os.Exit(1)
}

func output(w io.Writer, data interface{}, format string,
	outputTable func(w io.Writer, data interface{}) error,
) error {
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
	if _, err = w.Write(bytes); err != nil {
		return err
	}
	return nil
}

func outputJSON(w io.Writer, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	if _, err = w.Write(bytes); err != nil {
		return err
	}
	return nil
}
