package plugin

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var FixCmd = &cobra.Command{
	Use:   "fix",
	Short: "Clean up invalid enabled tenants in the plugin",
	Example: `
# Clean up all plugins
tkeel plugin fix

# Clean up the specified plugins
tkeel plugin fix [plugin-id ...]
`,
	Run: func(cmd *cobra.Command, args []string) {
		plugins := args
		print.PendingStatusEvent(os.Stdout, "Clean invalid enabled tenants....")
		daprStatus, err := kubernetes.CheckDapr()
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		tenantList, err := kubernetes.TenantList()
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		var tenants = make([]string, len(tenantList))
		for _, tenant := range tenantList {
			tenants = append(tenants, tenant.ID)
		}

		if len(plugins) == 0 {
			pluginList, err := kubernetes.InstalledPlugin()
			if err != nil {
				print.FailureStatusEvent(os.Stdout, err.Error())
				os.Exit(1)
			}
			plugins = make([]string, len(pluginList))
			for i, plugin := range pluginList {
				plugins[i] = plugin.Name
			}
		}

		for _, plugin := range plugins {
			err := kubernetes.CleanInvalidTenants(plugin, tenants, daprStatus.Namespace)
			if err != nil {
				print.FailureStatusEvent(os.Stdout, fmt.Sprintf("clean invalid tenants failed, plugin: %s, err: %s", plugin, err.Error()))
			}
		}
		print.InfoStatusEvent(os.Stdout, "Invalid tenant cleanup completed")
	},
}

func init() {
	FixCmd.Flags().BoolP("help", "h", false, "Print this help message")
	PluginCmd.AddCommand(FixCmd)
}
