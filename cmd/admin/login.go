package admin

import (
	"os"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
	"github.com/tkeel-io/cli/pkg/print"
	"golang.org/x/term"
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
		var (
			token string
			err   error
			bytes []byte
		)
		if len(args) == 0 && password == "" {
			bytes, err = term.ReadPassword(syscall.Stdin)
			if err != nil {
				print.FailureStatusEvent(os.Stdout, "failed to read password from stdin")
				return
			}
			password = string(bytes)
		}

		token, err = kubernetes.AdminLogin(password)
		if err != nil {
			print.FailureStatusEvent(os.Stdout, "Login Failed: %s", err.Error())
			return
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
