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

package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// deployAppCmd represents the deployApp command
var deployAppCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy an application",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		appDeployer := GetAppDeployer()
		app := GetApp(appID)
		operation, err := appDeployer.DeployApp(app)
		if err != nil {
			fmt.Fprintf(os.Stderr, "DeployApp error: %s\n", err)
		}
		fmt.Printf("operation = %+v\n", operation)
		err = WaitForCompletion(ctx, operation)
		if err != nil {
			fmt.Fprintf(os.Stderr, "WaitForCompletion error: %s\n", err)
		}
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
	deployAppCmd.Flags().StringVarP(&appID, "app-id", "a", "", "app-id for new Marathon app")
}