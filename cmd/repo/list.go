package repo

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/fmtutil"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "show the tkeel repository list",
	Run: func(cmd *cobra.Command, args []string) {
		list, err := kubernetes.ListRepo()
		if err != nil {
			print.FailureStatusEvent(os.Stdout, "unable to list tkeel repositories")
			os.Exit(1)
		}
		fmtutil.Output(list)
	},
}
