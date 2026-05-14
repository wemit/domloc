package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "dev"

func init() {
	rootCmd.Version = version
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version)
		},
	})
}
