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

	"github.com/spf13/cobra"
)

// ConfigCmd is a command from removing a tKeel installation.
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Print default install config.",
	Example: `
# Print default install config
tkeel config

# Save default install config to file
tkeel config > config.yaml
`,
	PreRun: func(cmd *cobra.Command, args []string) {
	},
	Run: func(cmd *cobra.Command, args []string) {
		config := `host:
  admin: admin.tkeel.io
  tenant: tkeel.io
middleware:
  cache:
    customized: false
    url: redis://:Biz0P8Xoup@tkeel-middleware-redis-master:6379/0
  database:
    customized: false
    url: mysql://root:a3fks=ixmeb82a@tkeel-middleware-mysql:3306/tkeelauth
  queue:
    customized: false
    url: kafka://tkeel-middleware-kafka-headless:9092
  search:
    customized: false
    url: elasticsearch://admin:admin@tkeel-middleware-elasticsearch-master:9200
  service_registry:
    customized: false
    url: etcd://tkeel-middleware-etcd:2379
port: "30080"
repo:
  name: tkeel
  url: https://tkeel-io.github.io/helm-charts
`

		fmt.Println(config)
	},
}

func init() {
	ConfigCmd.Flags().BoolP("help", "h", false, "Print this help message")
	RootCmd.AddCommand(ConfigCmd)
}
