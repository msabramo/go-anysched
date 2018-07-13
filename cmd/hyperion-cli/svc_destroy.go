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

	"github.com/spf13/cobra"
)

var (
	svcID string // svc ID of service we are going to destroy
)

// svcDestroyCmd represents the "svc destroy" command
var svcDestroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy a service",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		operation, err := getManager().DestroySvc(svcID)
		if err != nil {
			fmt.Fprintf(os.Stderr, "DestroySvc error: %s\n", err)
			os.Exit(1)
		}
		if operation != nil {
			_, err = operation.Wait(ctx)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error: %s\n", err)
				os.Exit(1)
			}
		}
		fmt.Printf("Service %q deleted.\n", svcID)
		os.Exit(0)
	},
}

func init() {
	svcCmd.AddCommand(svcDestroyCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// destroyAppCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	svcDestroyCmd.Flags().StringVarP(&svcID, "svc-id", "s", "", "svc-id of service to destroy")

}
