// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation and Dapr Contributors.
// Licensed under the MIT License.
// ------------------------------------------------------------
package installer

import (
	"github.com/spf13/cobra"
)

const UserHelpExample = `
# Manage installer.
tkeel installer show <repo>/<installerId>@v<version>
tkeel installer list -r <repo>
tkeel installer list -a
`

var InstallerCmd = &cobra.Command{
	Use:     "installer",
	Short:   "installer manager.",
	Example: UserHelpExample,
	Run: func(cmd *cobra.Command, args []string) {
		// Prompt help information If there is no parameter
		if len(args) == 0 {
			cmd.Help()
			return
		}
	},
}

func init() {
	//InstallerCmd.Flags().BoolVarP(&kubernetesMode, "kubernetes", "k", true, "List tenant's enabled plugins in a Kubernetes cluster")
	InstallerCmd.Flags().BoolP("help", "h", false, "Print this help message")
}
