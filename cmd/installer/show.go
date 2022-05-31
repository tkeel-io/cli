// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation and Dapr Contributors.
// Licensed under the MIT License.
// ------------------------------------------------------------

package installer

import (
	"os"

	"github.com/gocarina/gocsv"
	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/fmtutil"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
	"github.com/tkeel-io/cli/pkg/utils"
)

var InstallerInfoCmd = &cobra.Command{
	Use:   "show",
	Short: "Show installer.",
	Example: `
# Show the specified installer
tkeel installer show <repo-name>/<installer-id>@v<version>
`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			print.WarningStatusEvent(os.Stdout, "Please specify installer info")
			print.WarningStatusEvent(os.Stdout, "For example, tkeel installer show <repo-name>/<installer-id>[@<version>]")
			os.Exit(1)
		}
		tkeelRepo, installer, version := utils.ParseInstallArg(args[0], officialRepo)
		data, err := kubernetes.InstallerInfo(tkeelRepo, installer, version)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		table, err := gocsv.MarshalString(data)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}
		fmtutil.PrintTable(table)
	},
}

func init() {
	InstallerInfoCmd.Flags().BoolP("help", "h", false, "Print this help message")
	InstallerInfoCmd.MarkFlagRequired("tenant")
	InstallerCmd.AddCommand(InstallerInfoCmd)
}
