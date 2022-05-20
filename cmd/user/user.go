package user

import (
	"github.com/spf13/cobra"
)

const UserHelpExample = `
# Manage plugins.
tkeel user create <username> <password> -t <tenantId>
tkeel user show <userId> -t <tenantId>
tkeel user delete <userId> -t <tenantId>
tkeel user list -t <tenantId>
`

var UserCmd = &cobra.Command{
	Use:     "user",
	Short:   "tenant user manage.",
	Example: UserHelpExample,
	Run: func(cmd *cobra.Command, args []string) {
		// Prompt help information If there is no parameter
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

func init() {
	// UserCmd.Flags().BoolVarP(&kubernetesMode, "kubernetes", "k", true, "List tenant's enabled plugins in a Kubernetes cluster")
	UserCmd.Flags().BoolP("help", "h", false, "Print this help message")
}
