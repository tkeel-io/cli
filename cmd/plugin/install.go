package plugin

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
	"github.com/tkeel-io/cli/pkg/utils"
	"github.com/tkeel-io/kit/log"
)

const officialRepo = "tkeel"

var (
	configFile string
)

var PluginInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the plugin which you want",
	Example: `
# Install the latest version
tkeel plugin install <repo-name>/<installer-id> <plugin-id>
# Install the specified version
tkeel plugin install <repo-name>/<installer-id>@<version> <plugin-id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			print.WarningStatusEvent(os.Stdout, "Please specify the installer info and plugin id")
			print.WarningStatusEvent(os.Stdout, "For example, tkeel plugin install <repo-name>/<installer-id>[@<version>] <plugin-id>")
			os.Exit(1)
		}
		var configb []byte
		var err error
		name := args[1]
		repo, plugin, version := utils.ParseInstallArg(args[0], officialRepo)
		if configFile != "" {
			configFile, err = filepath.Abs(configFile)
			if err != nil {
				print.FailureStatusEvent(os.Stdout, "unable to read config file")
				os.Exit(1)
			}
			configb, err = ioutil.ReadFile(configFile)
			if err != nil {
				print.FailureStatusEvent(os.Stdout, "unable to read config file")
				os.Exit(1)
			}
		}

		if err := kubernetes.Install(repo, plugin, version, name, configb); err != nil {
			log.Warn("install failed", err)
			print.FailureStatusEvent(os.Stdout, "Install %q failed, Because: %s", plugin, err.Error())
			os.Exit(1)
		}
		print.SuccessStatusEvent(os.Stdout, "Install %q success! It's named %q in k8s", plugin, name)
	},
}

func init() {
	PluginInstallCmd.Flags().StringVarP(&configFile, "config", "", "", "The plugin config file.")

	PluginCmd.AddCommand(PluginInstallCmd)
}
