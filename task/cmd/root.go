package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "main",
	Short: "CLI task manager",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "Some issue occurred")
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(addCmd, listCmd, doneCmd)
}