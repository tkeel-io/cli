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
	Use:     "uninstall",
	Short:   "uninstall the plugin which you want",
	Example: PluginHelpExample,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			print.PendingStatusEvent(os.Stdout, "please input the plugin name what you installed.")
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
	PluginUninstallCmd.Flags().BoolVarP(&debugMode, "debug", "", false, "The log mode")
	PluginCmd.AddCommand(PluginUninstallCmd)
}
