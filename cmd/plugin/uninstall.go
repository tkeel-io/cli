package plugin

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
	"github.com/tkeel-io/kit/log"
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
		if err := kubernetes.UninstallPlugin(pluginID); err != nil {
			log.Warn("remove the plugin failed", err)
			print.FailureStatusEvent(os.Stdout, "Try to remove installed plugin %q failed, Because: %s", strings.Join(args, ","), err.Error())
			os.Exit(1)
		}
		print.SuccessStatusEvent(os.Stdout, "Remove %q success!", strings.Join(args, ","))
	},
}

func init() {
	PluginCmd.AddCommand(PluginUninstallCmd)
}
