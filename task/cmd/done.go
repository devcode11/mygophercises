package cmd

import (
	"fmt"
	"os"
	"strconv"

	"main/store"
	"github.com/spf13/cobra"
)

var doneCmd = &cobra.Command{
	Use: "done",
	Short: "Marks task as completed",
	Run: func (cmd *cobra.Command, args []string) {
		if len(args) <= 0 {
			fmt.Fprintln(os.Stderr, "Usage: task done <task numbers...>")
		}
		for _, arg := range args {
			num, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s is not a number" , arg)
			}
			err = store.Done(num)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	},
}