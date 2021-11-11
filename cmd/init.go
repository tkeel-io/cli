// ------------------------------------------------------------
// Copyright 2021 The tKeel Contributors.
// Licensed under the Apache License.
// ------------------------------------------------------------

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var (
	kubernetesMode bool
	debugMode      bool
	wait           bool
	timeout        uint
	runtimeVersion string
	initNamespace  string
	enableMTLS     bool
	enableHA       bool
	values         []string
)

var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "Install tKeel platform on dapr.",
	PreRun: func(cmd *cobra.Command, args []string) {
	},
	Example: `
# Initialize Keel in Kubernetes
tkeel init 

# Initialize Keel in Kubernetes and wait for the installation to complete (default timeout is 300s/5m)
tkeel init --wait --timeout 600
`,
	Run: func(cmd *cobra.Command, args []string) {
		print.PendingStatusEvent(os.Stdout, "Making the jump to hyperspace...")
		if kubernetesMode {
			config := kubernetes.InitConfiguration{
				Namespace:  initNamespace,
				Version:    runtimeVersion,
				EnableMTLS: enableMTLS,
				EnableHA:   enableHA,
				Args:       values,
				Wait:       wait,
				Timeout:    timeout,
				DebugMode:  debugMode,
			}
			err := kubernetes.Init(config)
			if err != nil {
				print.FailureStatusEvent(os.Stdout, err.Error())
				os.Exit(1)
			}
			successEvent := fmt.Sprintf("Success! tKeel Platform has been installed to namespace %s. To verify, run `tkeel plugin list -k' in your terminal. To get started, go here: https://tkeel.io/keel-getting-started", config.Namespace)
			print.SuccessStatusEvent(os.Stdout, successEvent)
			return
		}
		print.FailureStatusEvent(os.Stdout, "Error! tKeel Platform should be in Kubernetes mode")
	},
}

func init() {
	InitCmd.Flags().BoolVarP(&kubernetesMode, "kubernetes", "k", true, "Deploy tKeel to a Kubernetes cluster")
	InitCmd.Flags().StringVarP(&runtimeVersion, "runtime-version", "", "latest", "The version of the tKeel Platform to install, for example: 1.0.0")
	InitCmd.Flags().StringVarP(&initNamespace,
		"namespace",      /* name */
		"n",              /*shorthand */
		"tkeel-platform", /* value */
		"The Kubernetes namespace to install tKeel in") /* usage */
	InitCmd.Flags().String("network", "", "The Docker network on which to deploy the tKeel Platform")
	InitCmd.Flags().BoolVarP(&wait, "wait", "", true, "Wait for Plugins initialization to complete")
	InitCmd.Flags().UintVarP(&timeout, "timeout", "", 300, "The wait timeout for the Kubernetes installation")
	InitCmd.Flags().BoolVarP(&debugMode, "debug", "", false, "The log mode")
	InitCmd.Flags().BoolP("help", "h", false, "Print this help message")
	RootCmd.AddCommand(InitCmd)
}
