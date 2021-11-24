package plugin

import (
	"context"
	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/helm"
	"github.com/tkeel-io/cli/pkg/print"
	"github.com/tkeel-io/kit/log"
	"os"
	"strings"
)

var PluginRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove the plugin which you want",
	Example: `
# Get status of tKeel plugins from Kubernetes
tkeel plugin list -k
tkeel plugin list --installable || -i
tkeel plugin delete -k pluginID
tkeel plugin register -k pluginID
`,
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
		return
	},
}

func init() {
	PluginRemoveCmd.Flags().BoolP("help", "h", false, "Print this help message")
	PluginCmd.AddCommand(PluginRemoveCmd)
}
