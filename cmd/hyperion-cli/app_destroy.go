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
	"time"

	"github.com/spf13/cobra"
)

var (
	appID           string // app ID of app we are going to destroy
	timeoutDuration = 15 * time.Second
)

// destroyAppCmd represents the destroyApp command
var destroyAppCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy an application",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		operation, err := Manager().DestroyApp(appID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return
		}
		if operation == nil {
			fmt.Printf("App %q deleted.\n", appID)
			return
		}
		_, err = operation.Wait(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %s\n", err)
			return
		}
	},
}

func init() {
	appCmd.AddCommand(destroyAppCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// destroyAppCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	destroyAppCmd.Flags().StringVarP(&appID, "app-id", "a", "", "app-id for new Marathon app")
	destroyAppCmd.Flags().DurationVarP(&timeoutDuration, "timeout", "t", timeoutDuration, "Max time to wait for deploy to complete")

}
