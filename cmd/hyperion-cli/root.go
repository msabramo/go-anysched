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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"git.corp.adobe.com/abramowi/hyperion"
	_ "git.corp.adobe.com/abramowi/hyperion/dockerswarm"
	_ "git.corp.adobe.com/abramowi/hyperion/kubernetes"
	_ "git.corp.adobe.com/abramowi/hyperion/marathon"
	_ "git.corp.adobe.com/abramowi/hyperion/nomad"
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
	if err := viper.BindPFlag("env", rootCmd.PersistentFlags().Lookup("env")); err != nil {
		panic(err)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		viper.AddConfigPath("$HOME/.config/hyperion-cli")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./etc")
		viper.SetConfigName("hyperion-cli")
	}

	viper.SetEnvPrefix("HYPERIONCLI")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func getManager() hyperion.Manager {
	managerConfig := getManagerConfig()
	manager, err := hyperion.NewManager(managerConfig)
	if err != nil {
		if _, err = fmt.Fprintf(os.Stderr, "error: %s\n", err); err != nil {
			panic(err)
		}
		os.Exit(1)
	}
	return manager
}

func getManagerConfig() hyperion.ManagerConfig {
	env := viper.GetString("env")
	if env == "" {
		die(`
			No env set. Set it with:
			  * --env option on command-line
			  * "env" setting in config file
			  * HYPERIONCLI_ENV environment variable`)
	}
	envRootKey := fmt.Sprintf("envs.%s", env)
	if viper.Get(envRootKey) == nil {
		die("env was %q but there was no %q in config file: %s", env, envRootKey, viper.ConfigFileUsed())
	}
	managerType := viper.GetString(envRootKey + ".type")
	managerAddress := viper.GetString(envRootKey + ".address")
	return hyperion.ManagerConfig{Type: managerType, Address: managerAddress}
}
