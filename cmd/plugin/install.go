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
	debugMode  bool
	wait       bool
	timeout    uint
	configFile string
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
		var configb []byte
		var err error
		name := args[1]
		repo, plugin, version := utils.ParseInstallArg(args[0], officialRepo)
		if configFile != "" {
			configFile, err = filepath.Abs(configFile)
			if err != nil {
				print.FailureStatusEvent(os.Stdout, "unable to read config file")
				return
			}
			configb, err = ioutil.ReadFile(configFile)
			if err != nil {
				print.FailureStatusEvent(os.Stdout, "unable to read config file")
				return
			}
		}

		if err := kubernetes.Install(repo, plugin, version, name, configb); err != nil {
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
	PluginInstallCmd.Flags().StringVarP(&configFile, "config", "", "", "The plugin config file.")

	PluginCmd.AddCommand(PluginInstallCmd)
}
