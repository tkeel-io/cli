package user

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var UserDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete user in tenant.",
	Example: `
# Delete the user of tenant by user id
tkeel user delete <user-id> -t <tenant-id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			print.WarningStatusEvent(os.Stdout, "Please specify the user id")
			print.WarningStatusEvent(os.Stdout, "For example, tkeel user delete <user-id> -t <tenant-id>")
			os.Exit(1)
		}
		userID := args[0]
		err := kubernetes.TenantUserDelete(tenant, userID)
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
