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
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var (
	repo string
)

var PluginStatusCmd = &cobra.Command{
	Use:     "list",
	Short:   "Show the health status of tKeel plugins. Supported platforms: Kubernetes",
	Example: PluginHelpExample,
	Run: func(cmd *cobra.Command, args []string) {
		if repo != "" {
			list, err := kubernetes.ListPluginsFromRepo(repo)
			if err != nil {
				if errors.Is(err, kubernetes.ErrInvalidToken) {
					print.FailureStatusEvent(os.Stdout, "please login!")
					return
				}
				print.FailureStatusEvent(os.Stdout, "unable to list plugins:%s", err.Error())
				return
			}
			outputList(list, len(list))
			return
		}

		status, err := kubernetes.InstalledList()
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		if len(status) == 0 {
			print.FailureStatusEvent(os.Stdout, "No status returned. Is tKeel plugins not install in your cluster?")
			os.Exit(1)
		}

		outputList(status, len(status))
	},
}

func init() {
	PluginStatusCmd.Flags().BoolVarP(&kubernetesMode, "kubernetes", "k", true, "List tenant's enabled plugins in a Kubernetes cluster")
	PluginStatusCmd.Flags().BoolP("help", "h", false, "Print this help message")
	PluginStatusCmd.Flags().StringVarP(&repo, "repo", "r", "", "Show the plugin list of this repository")
	PluginCmd.AddCommand(PluginStatusCmd)
}
