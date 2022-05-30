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

	"github.com/tkeel-io/cli/pkg/kubernetes"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/print"
)

var PluginInfoCmd = &cobra.Command{
	Use:   "show",
	Short: "Show plugin info.",
	Example: `
# Show information of the specified plugin
tkeel plugin show <plugin-id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		// TODO 显示详细信息。 启用的租户等
		if len(args) == 0 {
			print.PendingStatusEvent(os.Stdout, "PluginID not fount ...\n # Manager plugins. \n tkeel plugin register plugin-id")
			os.Exit(1)
		}

		pluginID := args[0]
		status, err := kubernetes.PluginInfo(pluginID)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		if len(status) == 0 {
			print.FailureStatusEvent(os.Stdout, "No status returned. Is tKeel plugins not install in your cluster?")
			os.Exit(1)
		}

		outputList(status, len(status))
	},
}

func init() {
	PluginCmd.AddCommand(PluginInfoCmd)
}
