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
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var PluginRemoveCmd = &cobra.Command{
	Use:   "unregister",
	Short: "Unregister plugins from tKeel. Supported platforms: Kubernetes",
	Example: `
# Get status of tKeel plugins from Kubernetes
tkeel plugin list -k
tkeel plugin list --installable || -i
tkeel plugin install https://tkeel-io.github.io/helm-charts/<pluginName> <pluginID>
tkeel plugin install https://tkeel-io.github.io/helm-charts/<pluginName>@v0.1.0 <pluginID>
tkeel plugin uninstall -k <pluginID>
tkeel plugin register -k <pluginID>
tkeel plugin remove <pluginID>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			print.PendingStatusEvent(os.Stdout, "PluginID not fount ...\n # Manager plugins. in Kubernetes mode \n tkeel plugin register -k pluginID")
			return
		}
		if kubernetesMode {
			pluginID := args[0]
			rmP, err := kubernetes.Unregister(pluginID)
			if err != nil {
				print.FailureStatusEvent(os.Stdout, err.Error())
				os.Exit(1)
			}
			print.InfoStatusEvent(os.Stdout, fmt.Sprintf("Unregister plugin: %s", rmP))
			print.SuccessStatusEvent(os.Stdout, fmt.Sprintf("Success! Plugin<%s> has been Registered to tKeel Platform . To verify, run `tkeel plugin list -k' in your terminal. ", pluginID))
		}
	},
}

func init() {
	PluginRemoveCmd.Flags().BoolVarP(&kubernetesMode, "kubernetes", "k", true, "List tenant's enabled plugins in a Kubernetes cluster")
	PluginRemoveCmd.Flags().BoolP("help", "h", false, "Print this help message")
	PluginCmd.AddCommand(PluginRemoveCmd)
}
