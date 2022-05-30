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
	Short: "Show the tkeel repository list",
	Example: `
# Show the tkeel repository list
tkeel repo list
`,
	Run: func(cmd *cobra.Command, args []string) {
		list, err := kubernetes.ListRepo()
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		fmtutil.Output(list)
	},
}
