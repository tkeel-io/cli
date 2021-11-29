package plugin

import (
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/print"
	"github.com/tkeel-io/kit/log"
)

var (
	debugMode bool
	wait      bool
	timeout   uint
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
tkeel plugin install https://tkeel-io.github.io/helm-charts/auth auth
tkeel plugin install https://tkeel-io.github.io/helm-charts/auth@v0.1.0 auth
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
		urls := strings.Split(plugin, "/")
		if len(urls) < 2 {
			print.PendingStatusEvent(os.Stdout, "please input the plugin which you want and the name you want")
			return
		}
		repo := strings.Join(urls[:len(urls)-1], "/")
		plugin = urls[len(urls)-1]

		config := kubernetes.InitConfiguration{
			Wait:      wait,
			Timeout:   timeout,
			DebugMode: debugMode,
		}

		if err := kubernetes.InstallPlugin(config, name, repo, plugin, version); err != nil {
			log.Warn("install failed", err)
			print.FailureStatusEvent(os.Stdout, "Install %q failed, Because: %s", plugin, err.Error())
			return
		}
		print.SuccessStatusEvent(os.Stdout, "Install %q success! It's named %q in k8s", plugin, name)
	},
}

func init() {
	PluginInstallCmd.Flags().BoolVarP(&kubernetesMode, "kubernetes", "k", true, "Deploy tKeel to a Kubernetes cluster")
	PluginInstallCmd.Flags().BoolVarP(&wait, "wait", "", true, "Wait for Plugins initialization to complete")
	PluginInstallCmd.Flags().UintVarP(&timeout, "timeout", "", 300, "The wait timeout for the Kubernetes installation")
	PluginInstallCmd.Flags().BoolVarP(&debugMode, "debug", "", false, "The log mode")
	PluginInstallCmd.Flags().BoolP("help", "h", false, "Print this help message")
	PluginCmd.AddCommand(PluginInstallCmd)
}
