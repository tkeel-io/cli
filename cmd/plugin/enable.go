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

var PluginEnableCmd = &cobra.Command{
	Use:   "enable",
	Short: "Enable plugins for tenant.",
	Example: `
# Enable the plugin for tenant
tkeel plugin enable <plugin-id> -t <tenant-id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			print.WarningStatusEvent(os.Stdout, "Please specify the plugin id.")
			print.WarningStatusEvent(os.Stdout, "For example, tkeel plugin enable <plugin-id> -t <tenant-id>.")
			os.Exit(1)
		}

		pluginID := args[0]
		err := kubernetes.EnablePlugin(pluginID, tenant)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		print.SuccessStatusEvent(os.Stdout, fmt.Sprintf("Success! Plugin<%s> has been enabled for tenant<%s>.", pluginID, tenant))
	},
}

func init() {
	PluginEnableCmd.Flags().StringVarP(&tenant, "tenant", "t", "", "tenant id")
	PluginEnableCmd.MarkFlagRequired("tenant")
	PluginCmd.AddCommand(PluginEnableCmd)
}
