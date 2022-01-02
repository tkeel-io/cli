package core

import (
	"github.com/spf13/cobra"
	"github.com/tkeel-io/cli/pkg/kubernetes"
)

var selector string
var search string
var watch bool

// getCmd represents the get command.
var getCmd = &cobra.Command{
	Use:   "get [entityid]",
	Short: "Display one or more entitties",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			if watch {
				kubernetes.CoreWatch(args[0])
			} else {
				kubernetes.CoreGet(args[0])
			}
		} else {
			kubernetes.CoreList(search, selector)
		}
	},
}

func init() {
	CoreCmd.AddCommand(getCmd)

	getCmd.Flags().StringVarP(&selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2)")
	getCmd.Flags().StringVarP(&search, "search", "s", "", "search keyword (e.g. -s abc)")
	getCmd.Flags().BoolVarP(&watch, "watch", "w", false, "After listing/getting the requested object, watch for changes.")
}
