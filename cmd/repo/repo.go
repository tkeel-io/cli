package repo

import "github.com/spf13/cobra"

var RepoCmd = &cobra.Command{
	Use:   "repo",
	Short: "show tkeel repo",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
		}
	},
}

func init() {
	RepoCmd.AddCommand(AddCmd, ListCmd, DeleteCmd)
}
