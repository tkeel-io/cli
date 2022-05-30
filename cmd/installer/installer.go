// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation and Dapr Contributors.
// Licensed under the MIT License.
// ------------------------------------------------------------
package installer

import (
	"os"

	"github.com/spf13/cobra"
)

const UserHelpExample = `
# Manage installer.

# Show the specified installer
tkeel installer show <repo-name>/<installer-id>@v<version>

# List the installer for the specified repo
tkeel installer list -r <repo-name>

# List the installers for all repositories
tkeel installer list -a
`

var InstallerCmd = &cobra.Command{
	Use:     "installer",
	Short:   "Installer manage.",
	Example: UserHelpExample,
	Run: func(cmd *cobra.Command, args []string) {
		// Prompt help information If there is no parameter
		if len(args) == 0 {
			cmd.Help()
			os.Exit(1)
		}
	},
}

func init() {
	InstallerCmd.Flags().BoolP("help", "h", false, "Print this help message")
}
