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
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"git.corp.adobe.com/abramowi/hyperion"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hyperion-cli",
	Short: "A command for demoing the hyperion library",
	Long: `A command that demos the hyperion library, allowing the user
to deploy services to Marathon, Kubernetes, etc.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hyperion-cli.yaml)")
	rootCmd.PersistentFlags().StringP("env", "e", "", "environment to target")
	viper.BindPFlag("env", rootCmd.PersistentFlags().Lookup("env"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".hyperionapp" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".hyperion-cli")
	}

	viper.SetEnvPrefix("HYPERIONCLI")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func getManager() hyperion.Manager {
	managerConfig := hyperion.ManagerConfig{Type: hyperion.ManagerTypeKubernetes}
	// or alternatively one of the following:
	//
	// managerConfig := hyperonlib.ManagerConfig{
	// 	Type:    hyperion.ManagerTypeMarathon,
	// 	Address: "http://127.0.0.1:8080",
	// }
	// managerConfig := hyperonlib.ManagerConfig{
	// 	Type:    hyperion.ManagerTypeDockerSwarm,
	// 	Address: "http://127.0.0.1:2377",
	// }
	// managerConfig := hyperonlib.ManagerConfig{
	// 	Type:    hyperion.ManagerTypeNomad
	// 	Address: "http://127.0.0.1:4646",
	// }

	env := viper.GetString("env")
	managerType := hyperion.ManagerType(viper.GetString(fmt.Sprintf("envs.%s.type", env)))
	managerAddress := viper.GetString(fmt.Sprintf("envs.%s.address", env))
	managerConfig = hyperion.ManagerConfig{Type: managerType, Address: managerAddress}

	manager, err := hyperion.NewManager(managerConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	return manager
}
