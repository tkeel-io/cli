package plugin

import (
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
	"github.com/tkeel-io/kit/log"
)

var (
	debugMode    bool
	wait         bool
	timeout      uint
	tkeelVersion string
	secret       string
)

var PluginInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the plugin which you want",
	Example: PluginCmd.Example,
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
			Version:   tkeelVersion,
			Wait:      wait,
			Timeout:   timeout,
			DebugMode: debugMode,
			Secret:    secret,
		}

		if err := kubernetes.InstallPlugin(config, repo, plugin, name, version); err != nil {
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
	PluginInstallCmd.Flags().StringVarP(&secret, "secret", "", "changeme", "The secret of the tKeel Platform to install, for example: dix9vng")
	PluginInstallCmd.Flags().StringVarP(&tkeelVersion, "tkeel_version", "", "0.2.0", "The plugin depened tkeel version.")
	PluginInstallCmd.Flags().BoolP("help", "h", false, "Print this help message")
	PluginCmd.AddCommand(PluginInstallCmd)
}
