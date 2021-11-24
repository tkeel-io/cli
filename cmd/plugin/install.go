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

var PluginInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the plugin which you want",
	Example: `
# Get status of tKeel plugins from Kubernetes
tkeel plugin list -k
tkeel plugin list --installable || -i
tkeel plugin delete -k pluginID
tkeel plugin register -k pluginID
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			print.PendingStatusEvent(os.Stdout, "please input the plugin which you want and the name you want")
			return
		}
		pluginFormInput, name := args[0], args[1]
		plugin := pluginFormInput
		version := "latest"
		if sp := strings.Split(pluginFormInput, "@"); len(sp) == 2 {
			plugin, version = sp[0], sp[1]
		}
		if err := helm.Install(context.Background(), name, plugin, version); err != nil {
			log.Warn("install failed", err)
			print.FailureStatusEvent(os.Stdout, "Install %q failed, Because: %s", plugin, err.Error())
			return
		}
		print.SuccessStatusEvent(os.Stdout, "Install %q success! It's named %q in k8s", plugin, name)
		return
	},
}

func init() {
	PluginInstallCmd.Flags().BoolP("help", "h", false, "Print this help message")
	PluginCmd.AddCommand(PluginInstallCmd)
}
