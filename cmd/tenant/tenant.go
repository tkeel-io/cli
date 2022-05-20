package tenant

import (
	"github.com/spf13/cobra"
)

var TenantHelpExample = `
# Manage tenants.
tkeel tenant create <tenantTitle>
tkeel tenant show <tenantId>
tkeel tenant delete <tenantId>
tkeel tenant list
`
var TenantCmd = &cobra.Command{
	Use:     "tenant",
	Short:   "manage tenants.",
	Example: TenantHelpExample,
	Run: func(cmd *cobra.Command, args []string) {
		// Prompt help information If there is no parameter
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

func init() {
	TenantCmd.Flags().BoolP("help", "h", false, "Print this help message")
	// ListCmd.Flags().StringVarP(&tenant, "tenant", "t", "", "Tenant ID")
	// ListCmd.MarkFlagRequired("tenant")
}
