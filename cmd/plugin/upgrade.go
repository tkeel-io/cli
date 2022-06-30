package plugin

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
	"github.com/tkeel-io/cli/pkg/utils"
)

var PluginUpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the plugin which you want",
	Example: `
# Upgrade the specified plugin to the latest version
tkeel plugin upgrade <repo-name>/<installer-id> <plugin-id>

# Upgrade the specified plugin to the specified version
tkeel plugin upgrade <repo-name>/<installer-id>@<version> <plugin-id>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			print.WarningStatusEvent(os.Stdout, "Please specify the installer info and plugin id")
			print.WarningStatusEvent(os.Stdout, "For example, tkeel plugin upgrade <repo-name>/<installer-id>[@<version>] <plugin-id>")
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

		if err := kubernetes.PluginUpgrade(repo, plugin, version, name, configb); err != nil {
			print.FailureStatusEvent(os.Stdout, "Upgrade %q failed, Because: %s", plugin, err.Error())
			os.Exit(1)
		}
		print.SuccessStatusEvent(os.Stdout, "Upgrade %q success! It's named %q in k8s", plugin, name)
	},
}

func init() {
	PluginUpgradeCmd.Flags().StringVarP(&configFile, "config", "", "", "The plugin config file.")

	PluginCmd.AddCommand(PluginUpgradeCmd)
}
