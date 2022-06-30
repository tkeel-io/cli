package tenant

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var username string
var password string
var remark string
var TenantCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create new tenant.",
	Example: `
# Create tenant
tkeel tenant create [tenant-space-name]
`,
	Run: func(cmd *cobra.Command, args []string) {
		var title string
		if len(args) == 0 {
			err := survey.AskOne(&survey.Input{Message: "What the tenant space name?"}, &title)
			if err != nil {
				print.FailureStatusEvent(os.Stdout, err.Error())
				os.Exit(1)
			}
		} else {
			title = args[0]
		}
		if username == "" {
			err := survey.AskOne(&survey.Input{Message: "What the tenant admin username?"}, &username)
			if err != nil {
				print.FailureStatusEvent(os.Stdout, err.Error())
				os.Exit(1)
			}
		}
		if password == "" {
			err := survey.AskOne(&survey.Password{Message: "What the tenant admin password?"}, &password)
			if err != nil {
				print.FailureStatusEvent(os.Stdout, err.Error())
				os.Exit(1)
			}
		}
		err := kubernetes.TenantCreate(title, remark, username, password)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, err.Error())
			os.Exit(1)
		}

		print.SuccessStatusEvent(os.Stdout, "Successfully created!")
	},
}

func init() {
	TenantCreateCmd.Flags().BoolP("help", "h", false, "Print this help message")
	TenantCreateCmd.Flags().StringVarP(&username, "username", "u", "", "username of tenant")
	TenantCreateCmd.Flags().StringVarP(&password, "password", "p", "", "password of tenant")
	TenantCreateCmd.Flags().StringVarP(&remark, "remark", "r", "", "remark of tenant")
	TenantCmd.AddCommand(TenantCreateCmd)
}
