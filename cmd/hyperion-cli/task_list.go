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
	"fmt"
	"io"
	"os"

	"git.corp.adobe.com/abramowi/hyperion"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listTasksCmd = &cobra.Command{
	Use:   "list",
	Short: "List running tasks",
	Run: func(cmd *cobra.Command, args []string) {
		manager := getManager()
		tasks, err := manager.AllTasks()
		if err != nil {
			fmt.Fprintf(os.Stderr, "task list: AllTasks error: %s\n", err)
			return
		}
		err = output(os.Stdout, tasks, viper.GetString("output_format"), outputTaskListTable)
		if err != nil {
			fmt.Fprintf(os.Stderr, "task list: output error: %s\n", err)
			return
		}
	},
}

func outputTaskListTable(w io.Writer, data interface{}) error {
	tasks := data.([]hyperion.TaskInfo)
	for _, task := range tasks {
		fmt.Fprintf(w, "%-40s %-16s %-16s %s\n", task.Name, task.HostIP, task.TaskIP, task.ReadyTime)
	}
	return nil
}

func init() {
	taskCmd.AddCommand(listTasksCmd)

	listTasksCmd.Flags().StringP("output-format", "f", "yaml", `output format: "table", "yaml", "json"`)
	viper.BindPFlag("output_format", listTasksCmd.Flags().Lookup("output-format"))
}
