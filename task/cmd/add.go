package cmd

import (
	"fmt"
	"os"

	"main/store"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add task to list",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) <= 0 {
			fmt.Fprintln(os.Stderr, "Usage: task add <tasks...>")
		}
		for _, arg := range args {
			store.Add(arg)
		}
	},
}
