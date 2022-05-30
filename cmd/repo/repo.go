package repo

import "github.com/spf13/cobra"

var RepoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Repository manage.",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	RepoCmd.AddCommand(AddCmd, ListCmd, DeleteCmd)
}
