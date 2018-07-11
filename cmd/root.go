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
	"time"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	hyperionlib "git.corp.adobe.com/abramowi/hyperion/lib"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hyperion",
	Short: "A command for demoing the hyperion library",
	Long: `A command that demos the hyperion library, allowing the user
to deploy apps to Marathon, Kubernetes, etc.`,
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hyperionapp.yaml)")
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
		viper.SetConfigName(".hyperionapp")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func GetAppDeployer() hyperionlib.AppDeployer {
	appDeployerConfig := hyperionlib.AppDeployerConfig{
		Type:    hyperionlib.AppDeployerTypeKubernetes,
		Address: "kubeconfig",
	}
	// or alternatively one of the following:
	//
	// appDeployerConfig := AppDeployerConfig{
	// 	Type:    hyperionlib.AppDeployerTypeMarathon,
	// 	Address: "http://127.0.0.1:8080",
	// }
	// appDeployerConfig := AppDeployerConfig{
	// 	Type:    hyperionlib.AppDeployerTypeDockerSwarm,
	// 	Address: "http://127.0.0.1:2377",
	// }
	// appDeployerConfig := AppDeployerConfig{
	// 	Type:    hyperionlib.AppDeployerTypeNomad
	// 	Address: "http://127.0.0.1:4646",
	// }

	appDeployer, err := hyperionlib.NewAppDeployer(appDeployerConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
	return appDeployer
}

func WaitForCompletion(ctx context.Context, operation hyperionlib.Operation) error {
	if asyncOperation, ok := operation.(hyperionlib.AsyncOperation); ok && asyncOperation != nil {
		return asyncOperation.Wait(ctx, 15*time.Second)
	}
	return nil
}
