package admin

import (
	"github.com/spf13/cobra"
)

var AdminCmd = &cobra.Command{
	Use:   "admin",
	Short: "Use admin control the tkeel.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
