package core

import (
	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
)

// applyCmd represents the apply command.
var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply a configuration to a entity by filename",
	Run: func(cmd *cobra.Command, args []string) {
		kubernetes.CoreApply(filenames)
	},
}

func init() {
	CoreCmd.AddCommand(applyCmd)

	applyCmd.Flags().StringSliceVarP(&filenames, "filename", "f", []string{""}, "filename")
}
