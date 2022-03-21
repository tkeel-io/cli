package tenant

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var TenantDeleteCmd = &cobra.Command{
	Use:     "delete",
	Short:   "delete tenant info.",
	Example: TenantHelpExample,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			print.FailureStatusEvent(os.Stdout, "please input tenant id")
			return
		}
		tenantID := args[0]
		err := kubernetes.TenantDelete(tenantID)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		print.SuccessStatusEvent(os.Stdout, "Successfully delete!")

	},
}

func init() {
	TenantDeleteCmd.Flags().BoolP("help", "h", false, "Print this help message")
	TenantCmd.AddCommand(TenantDeleteCmd)
}
