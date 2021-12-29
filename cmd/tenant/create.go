// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation and Dapr Contributors.
// Licensed under the MIT License.
// ------------------------------------------------------------

package tenant

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var TenantCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "create tenant . Supported platforms: Kubernetes",
	Example: `
# Manager tenant. in Kubernetes mode
tk tenant create -k tenantTitle
tk tenant list -k
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			print.PendingStatusEvent(os.Stdout, "tenantTitle not fount ...\n # auth tenant. in Kubernetes mode \n tk auth createtenant -k tenantTitle adminName adminPassword")
			return
		}
		if kubernetesMode {
			title := args[0]
			adminName, adminPw := "", ""
			if len(args) == 3 {
				adminName = args[1]
				adminPw = args[2]
			}
			err := kubernetes.TenantCreate(title, adminName, adminPw)
			if err != nil {
				print.FailureStatusEvent(os.Stdout, err.Error())
				os.Exit(1)
			}

			print.SuccessStatusEvent(os.Stdout, "Success! ")
		}
	},
}

func init() {
	TenantCreateCmd.Flags().BoolVarP(&kubernetesMode, "kubernetes", "k", true, "List tenant's enabled plugins in a Kubernetes cluster")
	TenantCreateCmd.Flags().BoolP("help", "h", false, "Print this help message")
	TenantCreateCmd.MarkFlagRequired("kubernetes")
	TenantCmd.AddCommand(TenantCreateCmd)
}
