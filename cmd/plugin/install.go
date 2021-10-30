// ------------------------------------------------------------
// Copyright 2021 The TKeel Contributors.
// Licensed under the Apache License.
// ------------------------------------------------------------

package plugin

import (
	"os"

	"github.com/dapr/cli/utils"
	"github.com/gocarina/gocsv"
	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var PluginInstallCmd = &cobra.Command{
	Use:   "list",
	Short: "Show the health status of TKeel plugins. Supported platforms: Kubernetes",
	Example: `
# Get status of TKeel plugins from Kubernetes
tkeel plugin list -k
tkeel plugin delete -k pluginID
tkeel plugin register -k pluginID
`,
	Run: func(cmd *cobra.Command, args []string) {
		sc, err := kubernetes.NewStatusClient()
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		status, err := sc.Status()
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		if len(status) == 0 {
			print.FailureStatusEvent(os.Stdout, "No status returned. Is TKeel initialized in your cluster?")
			os.Exit(1)
		}
		table, err := gocsv.MarshalString(status)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}

		utils.PrintTable(table)
	},
}

func init() {
	PluginInstallCmd.Flags().BoolVarP(&kubernetesMode, "kubernetes", "k", true, "List tenant's enabled plugins in a Kubernetes cluster")
	PluginInstallCmd.Flags().BoolP("help", "h", false, "Print this help message")
	PluginCmd.AddCommand(PluginInstallCmd)
}
