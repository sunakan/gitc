package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gitc",
	Short: "Git repository cleanup tool",
	Long: `gitc is a CLI tool that automates Git repository cleanup.
It switches to the default branch, pulls the latest changes,
and removes unnecessary local branches.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("gitc - Git repository cleanup tool")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}