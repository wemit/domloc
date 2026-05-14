package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "domloc",
	Short:        "Local domain routing for developers",
	Long:         `Map domains to localhost ports. Zero config. Reliable HTTPS.`,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(
		initCmd(),
		addCmd(),
		removeCmd(),
		listCmd(),
		doctorCmd(),
		wildcardCmd(),
		resetCmd(),
		openCmd(),
	)
}
