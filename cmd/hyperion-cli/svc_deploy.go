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
	"github.com/spf13/viper"

	"git.corp.adobe.com/abramowi/hyperion"
)

var (
	deploySettings  = struct{ svcCfg hyperion.SvcCfg }{}
	timeoutDuration = 15 * time.Second
)

// svcDeployCmd represents the "svc deploy" command
var svcDeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy a service",
	Run: func(cmd *cobra.Command, args []string) {
		timeout := 60 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		startTime := time.Now()
		manager := getManager()
		deployment, err := manager.DeploySvc(deploySettings.svcCfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "DeploySvc error: %s\n", err)
			os.Exit(1)
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
				fmt.Fprintf(os.Stderr, "Deployment polling aborted after %v: %s\n", timeout, ctx.Err())
				return
			case <-time.After(1 * time.Second):
				status, err := deployment.GetStatus()
				if err != nil {
					fmt.Fprintf(os.Stderr, "GetStatus error: %s\n", err)
					if strings.Contains(err.Error(), "Not implemented") {
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
					tasks, err := manager.SvcTasks(deploySettings.svcCfg)
					if err != nil {
						fmt.Fprintf(os.Stderr, "app deploy: SvcTasks error: %s\n", err)
						if strings.Contains(err.Error(), "Not implemented") {
							return
						}
						continue
					}
					fmt.Printf("Deployment completed in %s\n\n", elapsedTime)
					err = output(os.Stdout, tasks, viper.GetString("output_format"), outputTaskListTable)
					if err != nil {
						fmt.Fprintf(os.Stderr, "app list: task list output error: %s\n", err)
						os.Exit(1)
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
	svcCmd.AddCommand(svcDeployCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deployAppCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	svcDeployCmd.Flags().StringVarP(&deploySettings.svcCfg.ID, "svc-id", "s", "", "ID for new service")
	svcDeployCmd.Flags().StringVarP(&deploySettings.svcCfg.Image, "image", "i", "", "Docker image for new service")
	svcDeployCmd.Flags().IntVarP(&deploySettings.svcCfg.Count, "count", "c", 1, "Number of containers to run")
	svcDeployCmd.Flags().DurationVarP(&timeoutDuration, "timeout", "t", timeoutDuration, "Max time to wait for deploy to complete")
}
