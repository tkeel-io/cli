// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation and Dapr Contributors.
// Licensed under the MIT License.
// ------------------------------------------------------------

package tenant

import (
	"github.com/gocarina/gocsv"
	"os"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
	"github.com/tkeel-io/cli/utils"
)

var TenantCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create tenant . Supported platforms: Kubernetes",
	Example: `
# Manager plugins. in Kubernetes mode
tKeel tenant create -k tenantTitle
tKeel tenant list -k
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			print.PendingStatusEvent(os.Stdout, "tenantTitle not fount ...\n # auth plugins. in Kubernetes mode \n tkeel auth createtenant -k tenantTitle")
			return
		}
		if kubernetesMode {
			title := args[0]
			data, err := kubernetes.TenantCreate(title)
			if err != nil {
				print.FailureStatusEvent(os.Stdout, err.Error())
				os.Exit(1)
			}
			dataSlice := []kubernetes.TenantCreateResp{*data}
			table, err := gocsv.MarshalString(dataSlice)
			if err != nil {
				print.FailureStatusEvent(os.Stdout, err.Error())
				os.Exit(1)
			}

			utils.PrintTable(table)
		}
	},
}

func init() {
	TenantCreateCmd.Flags().BoolVarP(&kubernetesMode, "kubernetes", "k", true, "List tenant's enabled plugins in a Kubernetes cluster")
	TenantCreateCmd.Flags().BoolP("help", "h", false, "Print this help message")
	TenantCreateCmd.MarkFlagRequired("kubernetes")
	TenantCmd.AddCommand(TenantCreateCmd)
}
