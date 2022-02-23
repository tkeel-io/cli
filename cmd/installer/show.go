// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation and Dapr Contributors.
// Licensed under the MIT License.
// ------------------------------------------------------------

package installer

import (
	"os"
	"strings"

	"github.com/gocarina/gocsv"
	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/fmtutil"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

const officialRepo = "tkeel"

var InstallerInfoCmd = &cobra.Command{
	Use:     "show",
	Short:   "show installer.",
	Example: UserHelpExample,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			print.FailureStatusEvent(os.Stdout, "please input installer info")
			return
		}
		repo, installer, version := parseInstallArg(args[0])
		data, err := kubernetes.InstallerInfo(repo, installer, version)
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

func parseInstallArg(arg string) (repo, plugin, version string) {
	version = ""
	plugin = arg

	if sp := strings.Split(arg, "@"); len(sp) == 2 {
		plugin, version = sp[0], sp[1]
	}

	if version != "" && version[0] == 'v' {
		version = version[1:]
	}

	repo = officialRepo
	if spi := strings.LastIndex(plugin, "/"); spi != -1 {
		repo, plugin = plugin[:spi], plugin[spi+1:]
		if repo == "" || strings.EqualFold(repo, "tkeel") {
			repo = officialRepo
			return
		}
	}
	return
}
