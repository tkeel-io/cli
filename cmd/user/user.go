package user

import (
	"github.com/spf13/cobra"
)

const UserHelpExample = `
# Manage plugins.
tkeel user create <username> <password> -t <tenant-id>
tkeel user show <user-id> -t <tenant-id>
tkeel user delete <user-id> -t <tenant-id>
tkeel user list -t <tenant-id>
`

var UserCmd = &cobra.Command{
	Use:     "user",
	Short:   "User manage of tenant.",
	Example: UserHelpExample,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	UserCmd.Flags().BoolP("help", "h", false, "Print this help message")
}
