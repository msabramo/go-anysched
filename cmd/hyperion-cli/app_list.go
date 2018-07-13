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

var listAppsCmd = &cobra.Command{
	Use:   "list",
	Short: "List running apps",
	Run: func(cmd *cobra.Command, args []string) {
		manager := getManager()
		apps, err := manager.AllApps()
		if err != nil {
			fmt.Fprintf(os.Stderr, "app list: AllApps error: %s\n", err)
			return
		}

		err = output(os.Stdout, apps, viper.GetString("output_format"), outputAppListTable)
		if err != nil {
			fmt.Fprintf(os.Stderr, "app list: output error: %s\n", err)
			return
		}
	},
}

func outputAppListTable(w io.Writer, data interface{}) error {
	apps := data.([]hyperion.AppInfo)
	for _, app := range apps {
		fmt.Fprintf(w, "%-40s\n", app.ID)
	}
	return nil
}

func init() {
	appCmd.AddCommand(listAppsCmd)

	listAppsCmd.Flags().StringP("output-format", "f", "yaml", `output format: "table", "yaml", "json"`)
	viper.BindPFlag("output_format", listTasksCmd.Flags().Lookup("output-format"))
}
