// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation and Dapr Contributors.
// Licensed under the MIT License.
// ------------------------------------------------------------

package user

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var UserDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "delete user of tenant.",
	Example: UserHelpExample,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			print.FailureStatusEvent(os.Stdout, "please input user id")
			return
		}
		userId := args[0]
		err := kubernetes.TenantUserDelete(tenant, userId)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		print.SuccessStatusEvent(os.Stdout, "Successfully delete!")
	},
}

func init() {
	UserDeleteCmd.Flags().BoolP("help", "h", false, "Print this help message")
	UserDeleteCmd.Flags().StringVarP(&tenant, "tenant", "t", "", "Tenant ID")
	UserDeleteCmd.MarkFlagRequired("tenant")
	UserCmd.AddCommand(UserDeleteCmd)
}
