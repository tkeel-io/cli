package plugin

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
	"github.com/tkeel-io/kit/log"
)

const officialRepo = "tkeel"

var (
	debugMode    bool
	wait         bool
	timeout      uint
	tkeelVersion string
	secret       string
	configFile   string
)

var PluginInstallCmd = &cobra.Command{
	Use:     "install",
	Short:   "Install the plugin which you want",
	Example: PluginHelpExample,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			cmd.Help()
			print.PendingStatusEvent(os.Stdout, "please input the plugin which you want and the name you want")
			return
		}
		var config string
		name := args[1]
		repo, plugin, version := parseInstallArg(args[0])
		if configFile != "" {
			configb, err := ioutil.ReadFile(configFile)
			if err != nil {
				print.FailureStatusEvent(os.Stdout, "unable to read config file")
				return
			}
			config = base64.StdEncoding.EncodeToString(configb)
		}

		if err := kubernetes.Install(repo, plugin, version, name, config); err != nil {
			log.Warn("install failed", err)
			print.FailureStatusEvent(os.Stdout, "Install %q failed, Because: %s", plugin, err.Error())
			return
		}
		print.SuccessStatusEvent(os.Stdout, "Install %q success! It's named %q in k8s", plugin, name)
	},
}

func init() {
	PluginInstallCmd.Flags().BoolVarP(&wait, "wait", "", true, "Wait for Plugins initialization to complete")
	PluginInstallCmd.Flags().UintVarP(&timeout, "timeout", "", 300, "The wait timeout for the Kubernetes installation")
	PluginInstallCmd.Flags().BoolVarP(&debugMode, "debug", "", false, "The log mode")
	PluginInstallCmd.Flags().StringVarP(&secret, "secret", "", "changeme", "The secret of the tKeel Platform to install, for example: dix9vng")
	PluginInstallCmd.Flags().StringVarP(&tkeelVersion, "tkeel_version", "", "0.2.0", "The plugin depened tkeel version.")
	PluginInstallCmd.Flags().StringVarP(&configFile, "config", "", "", "The plugin config file.")

	PluginCmd.AddCommand(PluginInstallCmd)
}

// parseInstallArg parse the first arg, get repo, plugin and version information.
// More efficient and concise support for both formats：
// url style install target plugin: https://tkeel-io.github.io/helm-charts/A@version
// short style install official plugin： tkeel/B@version or C@version.
func parseInstallArg(arg string) (repo, plugin, version string) {
	version = "latest"
	plugin = arg

	if sp := strings.Split(arg, "@"); len(sp) == 2 {
		plugin, version = sp[0], sp[1]
	}

	if version[0] == 'v' {
		version = version[1:]
	}

	repo = officialRepo
	if spi := strings.LastIndex(plugin, "/"); spi != -1 {
		repo, plugin = plugin[:spi], plugin[spi+1:]
		if repo == "" || strings.EqualFold(repo, "tkeel") {
			repo = officialRepo
			return
		}
	}
	return
}
