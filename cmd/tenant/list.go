package tenant

import (
	"os"

	"github.com/gocarina/gocsv"
	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/fmtutil"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var pluginID string
var TenantListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all tenant.",
	Example: `
# List tenant
tkeel tenant list
tkeel tenant list -p <pluginID>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if pluginID != "" {
			data, err := kubernetes.TenantPluginList(pluginID)
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
			os.Exit(1)
		}
		data, err := kubernetes.TenantList()
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
	TenantListCmd.Flags().BoolP("help", "h", false, "Print this help message")
	TenantListCmd.Flags().StringVarP(&pluginID, "plugin", "p", "", "List the tenant that enabled the plugin")
	TenantCmd.AddCommand(TenantListCmd)
}
