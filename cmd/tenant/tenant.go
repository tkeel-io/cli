// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation and Dapr Contributors.
// Licensed under the MIT License.
// ------------------------------------------------------------

package tenant

import (
	"github.com/spf13/cobra"
)

var TenantCmd = &cobra.Command{
	Use:   "tenant",
	Short: "tenant of Auth plugins. Supported platforms: Kubernetes",
	Example: `
# Manager plugins. in Kubernetes mode
tKeel tenant create -k tenantTitle
tKeel tenant list -k
`,
	Run: func(cmd *cobra.Command, args []string) {
		// Prompt help information If there is no parameter
		if len(args) == 0 {
			cmd.Help()
			return
		}
	},
}

func init() {
	TenantCmd.Flags().BoolVarP(&kubernetesMode, "kubernetes", "k", true, "List tenant's enabled plugins in a Kubernetes cluster")
	TenantCmd.Flags().BoolP("help", "h", false, "Print this help message")
	//ListCmd.Flags().StringVarP(&tenant, "tenant", "t", "", "Tenant ID")
	//ListCmd.MarkFlagRequired("tenant")
	RootCmd.AddCommand(TenantCmd)
}
