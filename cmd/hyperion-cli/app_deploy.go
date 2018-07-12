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
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"git.corp.adobe.com/abramowi/hyperion"
)

var (
	deploySettings = struct{ app hyperion.App }{}
)

// deployAppCmd represents the deployApp command
var deployAppCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an application",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		startTime := time.Now()
		manager := getManager()
		deployment, err := manager.DeployApp(deploySettings.app)
		if err != nil {
			fmt.Fprintf(os.Stderr, "DeployApp error: %s\n", err)
			return
		}
		for key, val := range deployment.GetProperties() {
			if key == "" || val == "" {
				continue
			}
			fmt.Printf("%-30s : %v\n", key, val)
		}
		fmt.Println()

		var lastUpdateTime time.Time

		for {
			select {
			case <-ctx.Done():
				fmt.Fprintf(os.Stderr, "Deployment polling error: %s\n", ctx.Err())
				return
			case <-time.After(1 * time.Second):
				status, err := deployment.GetStatus()
				if err != nil {
					fmt.Fprintf(os.Stderr, "GetStatus error: %s\n", err)
					if strings.Contains(err.Error(), "Not implemented yet") {
						return
					}
					continue
				}
				if status.LastUpdateTime == lastUpdateTime {
					continue
				}
				fmt.Printf("[%s] %s\n", status.LastUpdateTime.Format(time.RFC3339), status.Msg)
				lastUpdateTime = status.LastUpdateTime
				if status.Done {
					elapsedTime := time.Since(startTime)
					fmt.Printf("Deployment completed in %s\n", elapsedTime)
					tasks, err := manager.AppTasks(deploySettings.app)
					if err != nil {
						fmt.Fprintf(os.Stderr, "app deploy: AppTasks error: %s\n", err)
						continue
					}
					for _, task := range tasks {
						fmt.Printf("%-40s %-16s %-16s %s\n", task.Name, task.HostIP, task.TaskIP, task.ReadyTime)
					}
					return
				}
			}
		}
		/*
			_, err = deployment.Wait(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Deployment failed: %s\n", err)
			}
		*/
	},
}

func init() {
	appCmd.AddCommand(deployAppCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployAppCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	deployAppCmd.Flags().StringVarP(&deploySettings.app.ID, "app-id", "a", "", "app-id for new app")
	deployAppCmd.Flags().StringVarP(&deploySettings.app.Image, "image", "i", "", "Docker image for new app")
	deployAppCmd.Flags().IntVarP(&deploySettings.app.Count, "count", "c", 1, "Number of containers to run")
	deployAppCmd.Flags().DurationVarP(&timeoutDuration, "timeout", "t", timeoutDuration, "Max time to wait for deploy to complete")
}
