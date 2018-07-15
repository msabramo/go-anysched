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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"git.corp.adobe.com/abramowi/hyperion"
)

// svcListCmd represents the "svc list" command
var svcListCmd = &cobra.Command{
	Use:   "list",
	Short: "List running services",
	Run: func(cmd *cobra.Command, args []string) {
		manager := getManager()
		svcs, err := manager.Svcs()
		if err != nil {
			_, err2 := fmt.Fprintf(os.Stderr, "svc list: Svcs error: %s\n", err)
			if err2 != nil {
				panic(err2)
			}
			return
		}

		err = output(os.Stdout, svcs, viper.GetString("output_format"), outputSvcListTable)
		if err != nil {
			_, err2 := fmt.Fprintf(os.Stderr, "svc list: output error: %s\n", err)
			if err2 != nil {
				panic(err2)
			}
			return
		}
	},
}

func outputSvcListTable(w io.Writer, data interface{}) error {
	svcs := data.([]hyperion.Svc)
	for _, svc := range svcs {
		if _, err := fmt.Fprintf(w, "%-40s\n", svc.ID); err != nil {
			panic(err)
		}
	}
	return nil
}

func init() {
	svcCmd.AddCommand(svcListCmd)

	svcListCmd.Flags().StringP("output-format", "f", "yaml", `output format: "table", "yaml", "json"`)
	if err := viper.BindPFlag("output_format", taskListCmd.Flags().Lookup("output-format")); err != nil {
		panic(err)
	}
}
