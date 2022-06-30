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

	"github.com/spf13/cobra"

	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "CheckDapr tkeel install status.",
	Example: `
# CheckDapr the tkeel installation environment
tkeel status
`,
	Run: func(cmd *cobra.Command, args []string) {
		print.PendingStatusEvent(os.Stdout, "Checking the Dapr runtime status...")
		status, err := kubernetes.CheckTKeel()
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		installed := fmt.Sprintf("tkeel installed: %v", status.Installed)
		if status.Installed {
			print.SuccessStatusEvent(os.Stdout, installed)
		} else {
			print.FailureStatusEvent(os.Stdout, installed)
		}

		version := fmt.Sprintf("tkeel version: %s", status.Version)
		if status.Version != "" {
			print.SuccessStatusEvent(os.Stdout, version)
		} else {
			print.FailureStatusEvent(os.Stdout, version)
		}

		namespace := fmt.Sprintf("tkeel namespace: %s", status.Namespace)
		if status.Namespace != "" {
			print.SuccessStatusEvent(os.Stdout, namespace)
		} else {
			print.FailureStatusEvent(os.Stdout, namespace)
		}
	},
}

func init() {
	StatusCmd.Flags().BoolP("help", "h", false, "Print this help message")
	RootCmd.AddCommand(StatusCmd)
}
