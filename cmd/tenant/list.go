// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation and Dapr Contributors.
// Licensed under the MIT License.
// ------------------------------------------------------------

package tenant

import (
	"os"

	"github.com/gocarina/gocsv"
	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/fmtutil"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var TenantListCmd = &cobra.Command{
	Use:   "list",
	Short: "list tenant . Supported platforms: Kubernetes",
	Example: `
# Manager tenant. in Kubernetes mode
tk tenant create -k tenantTitle
tk tenant list -k
`,
	Run: func(cmd *cobra.Command, args []string) {
		if kubernetesMode {

			data, err := kubernetes.TenantList()
			if err != nil {
				print.FailureStatusEvent(os.Stdout, err.Error())
				os.Exit(1)
			}
			table, err := gocsv.MarshalString(data)
			if err != nil {
				print.FailureStatusEvent(os.Stdout, err.Error())
				os.Exit(1)
			}

			fmtutil.PrintTable(table)
		}
	},
}

func init() {
	TenantListCmd.Flags().BoolVarP(&kubernetesMode, "kubernetes", "k", true, "List tenant's enabled plugins in a Kubernetes cluster")
	TenantListCmd.Flags().BoolP("help", "h", false, "Print this help message")
	TenantListCmd.MarkFlagRequired("kubernetes")
	TenantCmd.AddCommand(TenantListCmd)
}
