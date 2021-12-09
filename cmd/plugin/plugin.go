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

//
//          developer         +        paas manager         +     tantent manager
//        						|                             |
//           +------------+       |       +-----------+         |      +----------+
//           |            |       |       |           |         |      |          |
//           | developing |       |       | published |         |      | disabled |
//           |            |       |       |           |         |      |          |
//           +----+-------+       |       +---+-------+         |      +---+------+
//        		|               |   install |                 |          |
//        		|  ^            |           v   ^             |          | ^
//        		|  |            |               | uninstall   |          | |
//        		|  |            |       +-------+---+         |          | |
//      release |  |            |       |           |         |   enable | |
//        		|  | upgrade    |       | installed |         |          | | disable
//        		|  |            |       |           |         |          | |
//        		|  |            |       +---+-------+         |          | |
//        		|  |            |  register |                 |          | |
//        		v  |            |           v  ^              |          v |
//        		   |            |              | remove       |            |
//           +-------+----+       |       +------+----+         |      +-----+----+
//           |            |       |       |           |         |      |          |
//           |  release   |       |       |registered |         |      | enabled  |
//           |            |       |       |           |         |      |          |
//           +------------+       +       +-----------+         +      +----------+

package plugin

import (
	"github.com/spf13/cobra"
)

var PluginHelpExample = `
# Get status of tKeel plugins from Kubernetes
tkeel plugin list -k
tkeel plugin list --installable || -i
tkeel plugin install https://tkeel-io.github.io/helm-charts/<pluginName> <pluginID>
tkeel plugin install https://tkeel-io.github.io/helm-charts/<pluginName>@v0.1.0 <pluginID>
tkeel plugin uninstall -k <pluginID>
tkeel plugin register -k <pluginID>
tkeel plugin unregister -k <pluginID>
`

var PluginCmd = &cobra.Command{
	Use:     "plugin",
	Short:   "Manager plugins. Supported platforms: Kubernetes",
	Example: PluginHelpExample,
	Run: func(cmd *cobra.Command, args []string) {
		// Prompt help information If there is no parameter
		if len(args) == 0 {
			cmd.Help()
			return
		}
	},
}

func init() {
	PluginCmd.PersistentFlags().BoolVarP(&kubernetesMode, "kubernetes", "k", true, "List tenant's enabled plugins in a Kubernetes cluster")
	PluginCmd.PersistentFlags().StringVarP(&outputFormat, "output", "o", "", "The output format of the list. Valid values are: json, yaml, or table (default)")
	PluginCmd.PersistentFlags().BoolP("help", "h", false, "Print this help message")
}
