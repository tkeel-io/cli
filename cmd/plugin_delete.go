// ------------------------------------------------------------
// Copyright 2021 The TKeel Contributors.
// Licensed under the Apache License.
// ------------------------------------------------------------

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var PluginDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Register plugins. Supported platforms: Kubernetes",
	Example: `
# Manager plugins. in Kubernetes mode
tkeel plugin list -k
tkeel plugin delete -k pluginID
tkeel plugin register -k pluginID
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			print.PendingStatusEvent(os.Stdout, "PluginId not fount ...\n # Manager plugins. in Kubernetes mode \n tkeel plugin delete -k pluginID")
			return
		}
		if kubernetesMode {
			pluginID := args[0]
			err := kubernetes.Delete(pluginID)
			if err != nil {
				print.FailureStatusEvent(os.Stdout, err.Error())
				os.Exit(1)
			}
			print.SuccessStatusEvent(os.Stdout, fmt.Sprintf("Success! Plugin<%s> has been deleted from TKeel Platform . To verify, run `tkeel plugin list -k' in your terminal. ", pluginID))
		}
	},
}

func init() {
	PluginDeleteCmd.Flags().BoolVarP(&kubernetesMode, "kubernetes", "k", true, "List tenant's enabled plugins in a Kubernetes cluster")
	PluginDeleteCmd.Flags().BoolP("help", "h", false, "Print this help message")
	PluginCmd.AddCommand(PluginDeleteCmd)
}
