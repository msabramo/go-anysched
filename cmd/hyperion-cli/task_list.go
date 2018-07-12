// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"git.corp.adobe.com/abramowi/hyperion"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v2"
)

var (
	outputFormat string
)

var listTasksCmd = &cobra.Command{
	Use:   "list",
	Short: "List running tasks",
	Run: func(cmd *cobra.Command, args []string) {
		manager := Manager()
		tasks, err := manager.AllTasks()
		if err != nil {
			fmt.Fprintf(os.Stderr, "task list: AllTasks error: %s\n", err)
			return
		}
		err = output(os.Stdout, tasks, outputFormat)
		if err != nil {
			fmt.Fprintf(os.Stderr, "task list: output error: %s\n", err)
			return
		}
	},
}

func output(w io.Writer, data interface{}, format string) error {
	switch format {
	case "yaml":
		return outputYAML(w, data)
	case "json":
		return outputJSON(w, data)
	case "table":
		return outputTable(w, data.([]hyperion.TaskInfo))
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

func outputTable(w io.Writer, tasks []hyperion.TaskInfo) error {
	for _, task := range tasks {
		fmt.Fprintf(w, "%-40s %-16s %-16s %s\n", task.Name, task.HostIP, task.TaskIP, task.ReadyTime)
	}
	return nil
}

func init() {
	taskCmd.AddCommand(listTasksCmd)

	listTasksCmd.Flags().StringVarP(&outputFormat, "output-format", "f", "table",
		`output format: "table", "yaml", "json"`)
}
