package user

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var UserCreateCmd = &cobra.Command{
	Use:     "create",
	Short:   "create user info.",
	Example: UserHelpExample,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			print.FailureStatusEvent(os.Stdout, "please input 1 arguments, 1st username, 2nd password")
			return
		}
		username := args[0]
		password := args[1]

		err := kubernetes.TenantUserCreate(tenant, username, password)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		print.SuccessStatusEvent(os.Stdout, "Success! ")
	},
}

func init() {
	// UserCreateCmd.Flags().BoolVarP(&kubernetesMode, "kubernetes", "k", true, "List tenant's enabled plugins in a Kubernetes cluster")
	UserCreateCmd.Flags().BoolP("help", "h", false, "Print this help message")
	UserCreateCmd.Flags().StringVarP(&tenant, "tenant", "t", "", "Tenant ID")
	UserCreateCmd.MarkFlagRequired("tenant")
	UserCmd.AddCommand(UserCreateCmd)
}
