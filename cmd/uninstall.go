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

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"

	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var (
	uninstallAll bool
	yes          bool
)

// UninstallCmd is a command from removing a tKeel installation.
var UninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall tKeel Platform.",
	Example: `
# Uninstall from kubernetes
tkeel uninstall
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		checkDapr()
	},
	Run: func(cmd *cobra.Command, args []string) {
		var confirm bool
		var err error
		if !yes {
			err := survey.AskOne(&survey.Confirm{Message: "Do you want to uninstall tkeel platform ?"}, &confirm)
			if err != nil {
				os.Exit(1)
			}
			if !confirm {
				os.Exit(0)
			}
		}

		defer func() {
			kubernetes.CleanToken()
		}()

		print.InfoStatusEvent(os.Stdout, "Removing tKeel Platform from your cluster...")

		if uninstallAll {
			err = kubernetes.UninstallAllPlugin(daprStatus.Namespace, debugMode)
			if err != nil {
				print.FailureStatusEvent(os.Stdout, fmt.Sprintf("Error removing plugins: %s", err))
			} else {
				print.SuccessStatusEvent(os.Stdout, "tKeel plugins has been removed successfully")
			}
		}

		err = kubernetes.UninstallPlatform(daprStatus.Namespace, timeout, debugMode)

		if err != nil {
			print.FailureStatusEvent(os.Stdout, fmt.Sprintf("Error removing tKeel: %s", err))
			os.Exit(1)
		} else {
			print.SuccessStatusEvent(os.Stdout, "tKeel Platform has been removed successfully")
		}
	},
}

func init() {
	UninstallCmd.Flags().UintVarP(&timeout, "timeout", "", 300, "The timeout for the Kubernetes uninstall")
	UninstallCmd.Flags().BoolVar(&uninstallAll, "all", false, "Remove @TODO .dapr directory, Redis, Placement and Zipkin containers")
	UninstallCmd.Flags().String("network", "", "The Docker network from which to remove the tKeel Platform")
	UninstallCmd.Flags().BoolVarP(&debugMode, "debug", "", false, "The log mode")
	UninstallCmd.Flags().BoolVarP(&yes, "yes", "y", false, "Uninstall the tkeel platform directly")
	UninstallCmd.Flags().BoolP("help", "h", false, "Print this help message")
	RootCmd.AddCommand(UninstallCmd)
}
