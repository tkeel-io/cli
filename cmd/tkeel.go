/*
Copyright 2021 The tKeel Authors.

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
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/tkeel-io/cli/cmd/plugin"
	"github.com/tkeel-io/cli/cmd/tenant"
	"github.com/tkeel-io/cli/pkg/api"
	"github.com/tkeel-io/cli/pkg/helm"
	"github.com/tkeel-io/cli/pkg/print"
)

var RootCmd = &cobra.Command{
	Use:   "tkeel",
	Short: "Keel CLI",
	Long: `
      __             __
     / /_____  ___  / /
    / //_/ _ \/ _ \/ /
   / ,< /  __/  __/ /
  /_/|_|\___/\___/_/
									   
===============================
Things Keel Platform`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if err := helm.SetNamespace(namespace); err != nil {
			fmt.Println("set helm namespace err")
			os.Exit(-1)
		}
	},
}

var (
	logAsJSON bool
	namespace string
)

// Execute adds all child commands to the root command.
func Execute(version, apiVersion string) {
	RootCmd.Version = version
	api.PlatformAPIVersion = apiVersion

	cobra.OnInitialize(initConfig)

	setVersion()

	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func setVersion() {
	template := fmt.Sprintf("Keel CLI version: %s \n", RootCmd.Version)
	RootCmd.SetVersionTemplate(template)
}

func initConfig() {
	if logAsJSON {
		print.EnableJSONFormat()
	}

	viper.SetEnvPrefix("keel")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&logAsJSON, "log-as-json", "", false, "Log output in JSON format")
	RootCmd.PersistentFlags().StringVarP(&namespace, "namespace", "n", "tkeel-platform", "The Kubernetes namespace to install tKeel in")
	RootCmd.AddCommand(plugin.PluginCmd)
	RootCmd.AddCommand(tenant.TenantCmd)
}
