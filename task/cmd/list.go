package cmd

import (
	"fmt"
	"main/store"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use: "list",
	Short: "List tasks to do",
	Run: func (cmd *cobra.Command, args []string) {
		tasks := store.List()
		if len(tasks) == 0 {
			fmt.Println("No tasks found")
			return
		}
		for i, task := range tasks {
			fmt.Printf("%d) %s\n", i+1, task)
		}
	},
}