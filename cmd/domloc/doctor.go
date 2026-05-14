package main

import (
	"github.com/spf13/cobra"
	"github.com/wemit/domloc/internal/doctor"
	"github.com/wemit/domloc/internal/ui"
)

func doctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Diagnose environment issues",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui.Header("Running diagnostics")
			checks := doctor.Run()
			doctor.Print(checks)
			return nil
		},
	}
}
