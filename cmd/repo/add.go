package repo

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add repository into tkeel",
	Example: `
# Add repository and specify a name
tkeel repo add <name> <url>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			print.WarningStatusEvent(os.Stdout, "Please specify the name and url.")
			print.WarningStatusEvent(os.Stdout, " For example, tkeel repo add <name> <url>")
			os.Exit(1)
		}
		name, url := args[0], args[1]
		err := kubernetes.AddRepo(name, url)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		print.SuccessStatusEvent(os.Stdout, "Successfully added!")
	},
}
