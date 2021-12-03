package plugin

import (
	"context"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/helm"
	"github.com/tkeel-io/cli/pkg/print"
	"github.com/tkeel-io/kit/log"
)

var PluginUninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "uninstall the plugin which you want",
	Example: PluginCmd.Example,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			print.PendingStatusEvent(os.Stdout, "please input the plugin name what you installed.")
			return
		}
		if err := helm.Uninstall(context.Background(), args...); err != nil {
			log.Warn("remove the plugin failed", err)
			print.FailureStatusEvent(os.Stdout, "Try to remove installed plugin %q failed, Because: %s", strings.Join(args, ","), err.Error())
			return
		}
		print.SuccessStatusEvent(os.Stdout, "Remove %q success!", strings.Join(args, ","))
	},
}

func init() {
	PluginUninstallCmd.Flags().BoolP("help", "h", false, "Print this help message")
	PluginCmd.AddCommand(PluginUninstallCmd)
}
