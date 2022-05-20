package admin

import (
	"os"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
)

var (
	password  string
	printable bool
)

var adminLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "login admin with password",
	Example: `tkeel admin login -p
	tkeel admin login
	tkeel admin login -p your_password --print
	`,
	Run: func(cmd *cobra.Command, args []string) {
		prompt := &survey.Password{Message: "Please enter your password: "}
		if len(args) == 0 && password == "" {
			if err := survey.AskOne(prompt, &password); err != nil {
				print.FailureStatusEvent(os.Stdout, "failed to read password from stdin")
				os.Exit(1)
			}
		}
		token, err := kubernetes.AdminLogin(password)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, "Login Failed: %s", err.Error())
			os.Exit(1)
		}

		print.SuccessStatusEvent(os.Stdout, "You are Login Success!")
		if printable {
			print.SuccessStatusEvent(os.Stdout, "Your Token: %s", token)
		}
	},
}

func init() {
	adminLoginCmd.Flags().StringVarP(&password, "password", "p", "", "input your admin password")
	adminLoginCmd.Flags().BoolVarP(&printable, "print", "", false, "print token after adminLoginCmd")
	AdminCmd.AddCommand(adminLoginCmd)
}
