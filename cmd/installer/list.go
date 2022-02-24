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
)

var InstallerListCmd = &cobra.Command{
	Use:     "list",
	Short:   "list installer.",
	Example: UserHelpExample,
	Run: func(cmd *cobra.Command, args []string) {
		if repo != "" {
			data, err := kubernetes.InstallerList(repo)
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
			return
		}
		if all {
			data, err := kubernetes.InstallerListAll()
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
			return
		}
		cmd.Help()
	},
}

func init() {
	InstallerListCmd.Flags().BoolP("help", "h", false, "Print this help message")
	InstallerListCmd.Flags().StringVarP(&repo, "repo", "r", "", "repo name")
	InstallerListCmd.Flags().BoolVarP(&all, "all-repo", "a", false, "all repo")
	InstallerCmd.AddCommand(InstallerListCmd)
}
