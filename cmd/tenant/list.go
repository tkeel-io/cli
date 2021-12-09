// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation and Dapr Contributors.
// Licensed under the MIT License.
// ------------------------------------------------------------

package tenant

import (
	"os"

	"github.com/tkeel-io/cli/fmtutil"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var TenantListCmd = &cobra.Command{
	Use:   "list",
	Short: "list tenant . Supported platforms: Kubernetes",
	Example: `
# Manager plugins. in Kubernetes mode
tKeel tenant create -k tenantTitle
tKeel tenant list -k
`,
	Run: func(cmd *cobra.Command, args []string) {
		if kubernetesMode {

			tenants, err := kubernetes.TenantList()
			if err != nil {
				print.FailureStatusEvent(os.Stdout, err.Error())
				os.Exit(1)
			}

			fmtutil.OutputList(tenants, len(tenants), outputFormat)
		}
	},
}

func init() {
	TenantListCmd.Flags().StringVarP(&outputFormat, "output", "o", "", "The output format of the list. Valid values are: json, yaml, or table (default)")
	TenantCmd.AddCommand(TenantListCmd)
}
