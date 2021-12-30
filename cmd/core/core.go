package core

import (
	"fmt"

	"github.com/spf13/cobra"
)

var filenames = make([]string, 0)

// CoreCmd represents the core command.
var CoreCmd = &cobra.Command{
	Use:   "core",
	Short: "core operation",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("core called")
	},
}

func init() {
}
