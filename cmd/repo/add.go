package repo

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "add repository into tkeel",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			print.FailureStatusEvent(os.Stdout, "please input 2 arguments,1st repo name 2nd repo url")
			os.Exit(1)
		}
		name, url := args[0], args[1]
		err := kubernetes.AddRepo(name, url)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, "unable add repo to tkeel")
			os.Exit(1)
		}
		print.SuccessStatusEvent(os.Stdout, "Successfully added!")
	},
}
