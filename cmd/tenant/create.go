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
# Manager plugins. in Kubernetes mode
tKeel tenant create tenantTitle
tKeel tenant list
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			print.PendingStatusEvent(os.Stdout, "tenantTitle not fount ...\n # auth plugins. in Kubernetes mode \n tkeel auth createtenant tenantTitle adminName adminPassword")
			return
		}
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
	},
}

func init() {
	TenantCreateCmd.Flags().BoolP("help", "h", false, "Print this help message")
	TenantCmd.AddCommand(TenantCreateCmd)
}
