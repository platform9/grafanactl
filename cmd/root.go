/*
Copyright © 2019 Platform9 Systems

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/platform9/grafana-sync/pkg/client"
	"github.com/spf13/cobra"

	"github.com/grafana-tools/sdk"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "grafana-sync",
	Short: "Sync dashboards across grafana organizations",
	Long: `Grafana Sync is a tool that enables replication of dashboards across
multiple grafana instances, or organizations.

You can download dashboards for a specific org, or folder.

You can upload dashboards to a specific org, preserving folder structure.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", `config file - default in order of precedence:
- .grafana-sync.yaml
- $HOME/.grafana-sync.yaml`)

	// These flags will override the config file if specified
	// `apiKey` command option for grafana API key
	rootCmd.PersistentFlags().String("apikey", "", "A Grafana API Key")
	viper.BindPFlag("apikey", rootCmd.PersistentFlags().Lookup("apikey"))
	// `url` command option for grafana URL
	rootCmd.PersistentFlags().String("url", "", "The URL of a Grafana instance")
	viper.BindPFlag("url", rootCmd.PersistentFlags().Lookup("url"))
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

		viper.SetConfigName(".grafana-sync.yaml")
		// First look in local directory
		viper.AddConfigPath(".")
		// Also look in HOME directory
		viper.AddConfigPath(home)
	}

	// Environment Variables expect to be the uppercase form of the flag name
	// env vars must be in the form GS_VARNAME
	viper.SetEnvPrefix("GS")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Println("Did not find config file. Continuing.")
		}
	}
}

// Ensures that the global authentication parameters are specified
// Will exit if they are not
func requireAuthParams() {
	if !viper.IsSet("url") {
		fmt.Fprintln(os.Stderr, "Error: Grafana URL not specified.")
		rootCmd.Println(rootCmd.UsageString())
		os.Exit(1)
	}
	if !viper.IsSet("apikey") {
		fmt.Fprintln(os.Stderr, "Error: Grafana APIKey not specified.")
		rootCmd.Println(rootCmd.UsageString())
		os.Exit(1)
	}
}

// Initializes a grafana client for the user
func getGrafanaClient() *sdk.Client {
	return sdk.NewClient(viper.GetString("url"), viper.GetString("apikey"), sdk.DefaultHTTPClient)
}

func getGrafanaClientInternal() *client.Client {
	return client.NewClient(viper.GetString("url"), viper.GetString("apikey"), client.DefaultHTTPClient)
}
