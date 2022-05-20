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

var PluginUpgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Short:   "upgrade the plugin which you want",
	Example: PluginHelpExample,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 2 {
			print.PendingStatusEvent(os.Stdout, "please input the plugin which you want and the name you want")
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
			log.Warn("upgrade failed", err)
			print.FailureStatusEvent(os.Stdout, "Upgrade %q failed, Because: %s", plugin, err.Error())
			os.Exit(1)
		}
		print.SuccessStatusEvent(os.Stdout, "Upgrade %q success! It's named %q in k8s", plugin, name)
	},
}

func init() {
	PluginUpgradeCmd.Flags().BoolVarP(&wait, "wait", "", true, "Wait for Plugins initialization to complete")
	PluginUpgradeCmd.Flags().UintVarP(&timeout, "timeout", "", 300, "The wait timeout for the Kubernetes installation")
	PluginUpgradeCmd.Flags().BoolVarP(&debugMode, "debug", "", false, "The log mode")
	PluginUpgradeCmd.Flags().StringVarP(&configFile, "config", "", "", "The plugin config file.")

	PluginCmd.AddCommand(PluginUpgradeCmd)
}
