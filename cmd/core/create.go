package core

import (
	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
)

// createCmd represents the create command.
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a entity from a file",
	Run: func(cmd *cobra.Command, args []string) {
		kubernetes.CoreCreate(filenames)
	},
}

func init() {
	CoreCmd.AddCommand(createCmd)

	createCmd.Flags().StringSliceVarP(&filenames, "filename", "f", []string{""}, "filename")
}
