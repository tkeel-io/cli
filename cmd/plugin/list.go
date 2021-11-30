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

package plugin

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/fmtutil"
	"github.com/tkeel-io/cli/pkg/helm"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
	"github.com/tkeel-io/kit/log"
)

var (
	installable bool
	update      bool
)

var PluginStatusCmd = &cobra.Command{
	Use:   "list",
	Short: "Show the health status of tKeel plugins. Supported platforms: Kubernetes",
	Example: `
# Get status of tKeel plugins from Kubernetes
tkeel plugin list -k
tkeel plugin list --installable || -i
tkeel plugin delete -k pluginID
tkeel plugin register -k pluginID
`,
	Run: func(cmd *cobra.Command, args []string) {
		if installable {
			if update {
				print.PendingStatusEvent(os.Stdout, "updating repo list")
			}
			list, err := helm.ListInstallable("table", update)
			if err != nil {
				log.Warn("list installable plugin failed.")
				print.FailureStatusEvent(os.Stdout, "list installable plugin failed. Because: %s", err.Error())
				return
			}
			fmt.Println(string(list))
			return
		}

		status, err := kubernetes.List()
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		if len(status) == 0 {
			print.FailureStatusEvent(os.Stdout, "No status returned. Is tKeel plugins not install in your cluster?")
			os.Exit(1)
		}
		csv, err := gocsv.MarshalString(status)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}

		fmtutil.PrintTable(csv)
	},
}

func init() {
	PluginStatusCmd.Flags().BoolVarP(&kubernetesMode, "kubernetes", "k", true, "List tenant's enabled plugins in a Kubernetes cluster")
	PluginStatusCmd.Flags().BoolP("help", "h", false, "Print this help message")
	PluginStatusCmd.Flags().BoolVarP(&update, "update", "", false, "this will update your repo list index")
	PluginStatusCmd.Flags().BoolVarP(&installable, "installable", "i", false, "Show the installable plugin")
	PluginCmd.AddCommand(PluginStatusCmd)
}
