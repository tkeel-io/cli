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

package plugin

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var PluginStatusCmd = &cobra.Command{
	Use:   "list",
	Short: "List all installed plugin in tkeel.",
	Example: `
# List the installed plugins 
tkeel plugin list
`,
	Run: func(cmd *cobra.Command, args []string) {
		if tenant != "" {
			list, err := kubernetes.ListPluginsFromTenant(tenant)
			if err != nil {
				print.FailureStatusEvent(os.Stdout, "unable to list plugins:%s", err.Error())
				os.Exit(1)
			}
			outputList(list, len(list))
			os.Exit(0)
		}

		status, err := kubernetes.InstalledList()
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		if len(status) == 0 {
			print.WarningStatusEvent(os.Stdout, "There is not plugin in your cluster.")
			os.Exit(0)
		}

		outputList(status, len(status))
	},
}

func init() {
	PluginStatusCmd.Flags().BoolP("help", "h", false, "Print this help message")
	//PluginStatusCmd.Flags().BoolVarP(&latest, "latest", "l", false, "Only show the latest plugin list of this repository")
	//PluginStatusCmd.Flags().StringVarP(&repo, "repo", "r", "", "Show the plugin list of this repository")
	PluginStatusCmd.Flags().StringVarP(&tenant, "tenant", "t", "", "Show the plugin of this tenant")
	PluginCmd.AddCommand(PluginStatusCmd)
}
