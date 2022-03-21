/*
Copyright 2021 The tKeel Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
	kitconfig "github.com/tkeel-io/kit/config"
)

var UpgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade tKeel platform.",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkDapr()
	},
	Example: `
# Initialize Keel in Kubernetes
tkeel init 

# Initialize Keel in Kubernetes and wait for the installation to complete (default timeout is 300s/5m)
tkeel init --wait --timeout 600
`,
	Run: func(cmd *cobra.Command, args []string) {
		print.PendingStatusEvent(os.Stdout, "Making the jump to hyperspace...")
		config := kubernetes.InitConfiguration{
			Namespace:  daprStatus.Namespace,
			Version:    runtimeVersion,
			EnableMTLS: enableMTLS,
			EnableHA:   enableHA,
			Args:       values,
			Wait:       wait,
			Timeout:    timeout,
			DebugMode:  debugMode,
			Repo: &kitconfig.Repo{
				Url:  repoURL,
				Name: repoName,
			},
		}
		err := kubernetes.Upgrade(config)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		successEvent := "Success! tKeel Platform upgrade success!"
		print.SuccessStatusEvent(os.Stdout, successEvent)
	},
}

func init() {
	UpgradeCmd.Flags().StringVarP(&runtimeVersion, "runtime-version", "", "latest", "The version of the tKeel Platform to install, for example: 1.0.0")
	UpgradeCmd.Flags().String("network", "", "The Docker network on which to deploy the tKeel Platform")
	UpgradeCmd.Flags().BoolVarP(&wait, "wait", "", true, "Wait for Plugins initialization to complete")
	UpgradeCmd.Flags().UintVarP(&timeout, "timeout", "", 300, "The wait timeout for the Kubernetes installation")
	UpgradeCmd.Flags().BoolVarP(&debugMode, "debug", "", false, "The log mode")
	UpgradeCmd.Flags().StringVarP(&configFile, "config", "f", "", "The tkeel installation config file")
	UpgradeCmd.Flags().StringVarP(&repoURL, "repo-url", "", "https://tkeel-io.github.io/helm-charts/", "The tkeel repo url")
	UpgradeCmd.Flags().StringVarP(&repoName, "repo-name", "", "tkeel", "The tkeel repo name")
	UpgradeCmd.Flags().BoolP("help", "h", false, "Print this help message")
	RootCmd.AddCommand(UpgradeCmd)
}
