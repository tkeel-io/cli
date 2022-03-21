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
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
	kitconfig "github.com/tkeel-io/kit/config"
)

var (
	debugMode      bool
	wait           bool
	timeout        uint
	runtimeVersion string
	secret         string
	enableMTLS     bool
	enableHA       bool
	values         []string
	configFile     string
	repoURL        string
	repoName       string
	password       string
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Install tKeel platform on dapr.",
	PreRun: func(cmd *cobra.Command, args []string) {
		checkDapr()
		if configFile != "" && configFile[0] == '~' {
			home, err := os.UserHomeDir()
			if err != nil {
				print.FailureStatusEvent(os.Stdout, err.Error())
				os.Exit(1)
			}
			configFile = strings.Replace(configFile, "~", home, 1)
		}
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
			Secret:     secret,
			Repo: &kitconfig.Repo{
				Url:  repoURL,
				Name: repoName,
			},
			Password:   password,
			ConfigFile: configFile,
		}
		err := kubernetes.Init(config)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		successEvent := fmt.Sprintf("Success! tKeel Platform has been installed to namespace %s. To verify, run `tkeel plugin list' in your terminal. To get started, go here: https://tkeel.io/keel-getting-started", config.Namespace)
		print.SuccessStatusEvent(os.Stdout, successEvent)
	},
}

func init() {
	InitCmd.Flags().StringVarP(&runtimeVersion, "runtime-version", "", "latest", "The version of the tKeel Platform to install, for example: 1.0.0")
	InitCmd.Flags().StringVarP(&secret, "secret", "", "changeme", "The secret of the tKeel Platform to install, for example: dix9vng")
	InitCmd.Flags().String("network", "", "The Docker network on which to deploy the tKeel Platform")
	InitCmd.Flags().BoolVarP(&wait, "wait", "", true, "Wait for Plugins initialization to complete")
	InitCmd.Flags().UintVarP(&timeout, "timeout", "", 300, "The wait timeout for the Kubernetes installation")
	InitCmd.Flags().BoolVarP(&debugMode, "debug", "", false, "The log mode")
	InitCmd.Flags().StringVarP(&configFile, "config", "f", "", "The tkeel installation config file")
	InitCmd.Flags().StringVarP(&repoURL, "repo-url", "", "https://tkeel-io.github.io/helm-charts/", "The tkeel repo url")
	InitCmd.Flags().StringVarP(&repoName, "repo-name", "", "tkeel", "The tkeel repo name")
	InitCmd.Flags().StringVarP(&password, "password", "", "changeme", "The tkeel admin password")
	InitCmd.Flags().BoolP("help", "h", false, "Print this help message")
	RootCmd.AddCommand(InitCmd)
}
