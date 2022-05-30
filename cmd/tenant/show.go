package tenant

import (
	"os"

	"github.com/gocarina/gocsv"
	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/fmtutil"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var TenantInfoCmd = &cobra.Command{
	Use:   "show",
	Short: "Show tenant info.",
	Example: `
# Show tenant info by tenant id
tkeel tenant show <tenant-id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			print.WarningStatusEvent(os.Stdout, "Please specify the tenant id")
			print.WarningStatusEvent(os.Stdout, "For example, tkeel tenant show <tenant-id>")
			os.Exit(1)
		}
		tenantID := args[0]
		data, err := kubernetes.TenantInfo(tenantID)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		table, err := gocsv.MarshalString(data)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		fmtutil.PrintTable(table)
	},
}

func init() {
	TenantInfoCmd.Flags().BoolP("help", "h", false, "Print this help message")
	TenantCmd.AddCommand(TenantInfoCmd)
}
