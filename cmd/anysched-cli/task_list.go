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

	"github.com/msabramo/go-anysched"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// taskListCmd represents the "task list" command
var taskListCmd = &cobra.Command{
	Use:   "list",
	Short: "List running tasks",
	Run: func(cmd *cobra.Command, args []string) {
		manager := getManager()
		tasks, err := manager.Tasks()
		if err != nil {
			_, err2 := fmt.Fprintf(os.Stderr, "task list: Tasks error: %s\n", err)
			if err2 != nil {
				panic(err2)
			}
			return
		}
		err = output(os.Stdout, tasks, viper.GetString("output_format"), outputTaskListTable)
		if err != nil {
			_, err2 := fmt.Fprintf(os.Stderr, "task list: output error: %s\n", err)
			if err2 != nil {
				panic(err2)
			}
			return
		}
	},
}

func outputTaskListTable(w io.Writer, data interface{}) error {
	tasks := data.([]anysched.Task)
	for _, task := range tasks {
		_, err := fmt.Fprintf(w, "%-40s %-16s %-16s %s\n", task.Name, task.HostIP, task.TaskIP, task.ReadyTime)
		if err != nil {
			panic(err)
		}
	}
	return nil
}

func init() {
	taskCmd.AddCommand(taskListCmd)

	taskListCmd.Flags().StringP("output-format", "f", "yaml", `output format: "table", "yaml", "json"`)
	if err := viper.BindPFlag("output_format", taskListCmd.Flags().Lookup("output-format")); err != nil {
		panic(err)
	}
}
