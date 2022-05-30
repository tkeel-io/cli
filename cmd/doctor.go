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

var DoctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Check tkeel install environment.",
	PreRun: func(cmd *cobra.Command, args []string) {
	},
	Example: `
# Check the tkeel installation environment
tkeel doctor
`,
	Run: func(cmd *cobra.Command, args []string) {
		print.PendingStatusEvent(os.Stdout, "Checking the Dapr runtime status...")
		status := kubernetes.Check()
		installed := fmt.Sprintf("dapr installed: %v", status.Installed)
		if status.Installed {
			print.SuccessStatusEvent(os.Stdout, installed)
		} else {
			print.FailureStatusEvent(os.Stdout, installed)
		}

		version := fmt.Sprintf("dapr version: %s", status.Version)
		if status.Version != "" {
			print.SuccessStatusEvent(os.Stdout, version)
		} else {
			print.FailureStatusEvent(os.Stdout, version)
		}

		namespace := fmt.Sprintf("dapr namespace: %s", status.Namespace)
		if status.Namespace != "" {
			print.SuccessStatusEvent(os.Stdout, namespace)
		} else {
			print.FailureStatusEvent(os.Stdout, namespace)
		}

		mtlsEnabled := fmt.Sprintf("dapr mtls enabled: %v", status.MTLSEnabled)
		if status.MTLSEnabled {
			print.SuccessStatusEvent(os.Stdout, mtlsEnabled)
		} else {
			print.FailureStatusEvent(os.Stdout, mtlsEnabled)
		}
	},
}

func init() {
	DoctorCmd.Flags().BoolP("help", "h", false, "Print this help message")
	RootCmd.AddCommand(DoctorCmd)
}
