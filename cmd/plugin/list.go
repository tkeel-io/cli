// ------------------------------------------------------------
// Copyright 2021 The tKeel Contributors.
// Licensed under the Apache License.
// ------------------------------------------------------------

package plugin

import (
	"os"

	"github.com/gocarina/gocsv"
	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/fmtutil"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var PluginStatusCmd = &cobra.Command{
	Use:   "list",
	Short: "Show the health status of tKeel plugins. Supported platforms: Kubernetes",
	Example: `
# Get status of tKeel plugins from Kubernetes
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
			print.FailureStatusEvent(os.Stdout, "No status returned. Is tKeel initialized in your cluster?")
			os.Exit(1)
		}
		csv, err := gocsv.MarshalString(status)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}

		fmtutil.PrintTable(csv)
	},
}

func init() {
	PluginStatusCmd.Flags().BoolVarP(&kubernetesMode, "kubernetes", "k", true, "List tenant's enabled plugins in a Kubernetes cluster")
	PluginStatusCmd.Flags().BoolP("help", "h", false, "Print this help message")
	PluginCmd.AddCommand(PluginStatusCmd)
}
