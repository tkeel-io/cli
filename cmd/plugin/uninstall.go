package plugin

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var PluginUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall the plugin which you want",
	Example: `
# Uninstall the specified plugin by id
tkeel plugin uninstall <plugin-id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			print.WarningStatusEvent(os.Stdout, "Please specify the plugin id.")
			print.WarningStatusEvent(os.Stdout, "For example, tkeel plugin uninstall <plugin-id>")
			os.Exit(1)
		}
		pluginID := args[0]
		if force {
			tenantList, err := kubernetes.TenantList()
			if err != nil {
				print.FailureStatusEvent(os.Stdout, "TenantList", err.Error())
				os.Exit(1)
			}
			for _, tenant := range tenantList {
				err = kubernetes.DisablePlugin(pluginID, tenant.ID)
				if err != nil {
					print.FailureStatusEvent(os.Stdout, "DisablePlugin, %s - %s. err: %s,", pluginID, tenant.ID, err.Error())
					os.Exit(1)
				}
			}
		}
		if err := kubernetes.UninstallPlugin(pluginID); err != nil {
			print.FailureStatusEvent(os.Stdout, "Try to remove installed plugin %q failed, Because: %s", strings.Join(args, ","), err.Error())
			os.Exit(1)
		}
		print.SuccessStatusEvent(os.Stdout, "Remove %q success!", strings.Join(args, ","))
	},
}

func init() {
	PluginUninstallCmd.Flags().BoolVarP(&force, "force", "f", false, "force uninstall plugin, even if a tenant has enabled it.")
	PluginCmd.AddCommand(PluginUninstallCmd)
}
