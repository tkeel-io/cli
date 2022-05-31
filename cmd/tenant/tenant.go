package tenant

import (
	"github.com/spf13/cobra"
)

var TenantHelpExample = `
# Manage tenants.
tkeel tenant create <tenant-name>
tkeel tenant show <tenant-id>
tkeel tenant delete <tenant-id>
tkeel tenant list
`
var TenantCmd = &cobra.Command{
	Use:     "tenant",
	Short:   "Tenant manage.",
	Example: TenantHelpExample,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	TenantCmd.Flags().BoolP("help", "h", false, "Print this help message")
}
