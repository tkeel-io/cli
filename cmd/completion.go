// ------------------------------------------------------------
// Copyright 2021 The TKeel Contributors.
// Licensed under the Apache License.
// ------------------------------------------------------------

package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionExample = `
	# Installing bash completion on macOS using homebrew
	## If running Bash 3.2 included with macOS
	brew install bash-completion
	## or, if running Bash 4.1+
	brew install bash-completion@2
	## Add the completion to your completion directory
	tkeel completion bash > $(brew --prefix)/etc/bash_completion.d/tkeel
	source ~/.bash_profile

	# Installing bash completion on Linux
	## If bash-completion is not installed on Linux, please install the 'bash-completion' package
	## via your distribution's package manager.
	## Load the tkeel completion code for bash into the current shell
	source <(tkeel completion bash)
	## Write bash completion code to a file and source if from .bash_profile
	tkeel completion bash > ~/.tkeel/completion.bash.inc
	printf "
	## tkeel shell completion
	source '$HOME/.tkeel/completion.bash.inc'
	" >> $HOME/.bash_profile
	source $HOME/.bash_profile

	# Installing zsh completion on macOS using homebrew
	## If zsh-completion is not installed on macOS, please install the 'zsh-completion' package
	brew install zsh-completions
	## Set the tkeel completion code for zsh[1] to autoload on startup
	tkeel completion zsh > "${fpath[1]}/_tkeel"
	source ~/.zshrc

	# Installing zsh completion on Linux
	## If zsh-completion is not installed on Linux, please install the 'zsh-completion' package
	## via your distribution's package manager.
	## Load the tkeel completion code for zsh into the current shell
	source <(tkeel completion zsh)
	# Set the tkeel completion code for zsh[1] to autoload on startup
	tkeel completion zsh > "${fpath[1]}/_tkeel"

	# Installing powershell completion on Windows
	## Create $PROFILE if it not exists
	if (!(Test-Path -Path $PROFILE )){ New-Item -Type File -Path $PROFILE -Force }
	## Add the completion to your profile
	tkeel completion powershell >> $PROFILE
`

func newCompletionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "completion",
		Short:   "Generates shell completion scripts",
		Example: completionExample,
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	cmd.AddCommand(
		newCompletionBashCmd(),
		newCompletionZshCmd(),
		newCompletionPowerShellCmd(),
	)

	cmd.Flags().BoolP("help", "h", false, "Print this help message")

	return cmd
}

func newCompletionBashCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bash",
		Short: "Generates bash completion scripts",
		Run: func(cmd *cobra.Command, args []string) {
			_ = RootCmd.GenBashCompletion(os.Stdout)
		},
	}

	cmd.Flags().BoolP("help", "h", false, "Print this help message")

	return cmd
}

func newCompletionZshCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zsh",
		Short: "Generates zsh completion scripts",
		Run: func(cmd *cobra.Command, args []string) {
			_ = RootCmd.GenZshCompletion(os.Stdout)
		},
	}
	cmd.Flags().BoolP("help", "h", false, "Print this help message")

	return cmd
}

func newCompletionPowerShellCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "powershell",
		Short: "Generates powershell completion scripts",
		Run: func(cmd *cobra.Command, args []string) {
			_ = RootCmd.GenPowerShellCompletion(os.Stdout)
		},
	}
	cmd.Flags().BoolP("help", "h", false, "Print this help message")

	return cmd
}

func init() {
	RootCmd.AddCommand(newCompletionCmd())
}
