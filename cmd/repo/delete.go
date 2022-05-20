package repo

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var DeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete the target tkeel repository",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			print.FailureStatusEvent(os.Stdout, "please input repo name which is you want to delete")
			os.Exit(1)
		}
		name := args[0]
		err := kubernetes.DeleteRepo(name)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, "unable delete repo to tkeel")
			os.Exit(1)
		}
		print.SuccessStatusEvent(os.Stdout, "Successfully delete!")
	},
}
