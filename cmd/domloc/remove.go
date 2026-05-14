package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/wemit/domloc/internal/caddy"
	"github.com/wemit/domloc/internal/registry"
	"github.com/wemit/domloc/internal/ui"
)

func removeCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "remove <domain>",
		Aliases: []string{"rm"},
		Short:   "Remove a domain route",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			domain := args[0]

			reg, err := registry.Load()
			if err != nil {
				return wrap("load registry", err)
			}

			if !reg.Remove(domain) {
				return fmt.Errorf("domain %q not found", domain)
			}

			if err := reg.Save(); err != nil {
				return wrap("save registry", err)
			}

			if err := caddy.EnsureRoutes(reg.Routes); err != nil {
				return wrap("reload Caddy", err)
			}

			ui.OK(fmt.Sprintf("removed %s", domain))
			return nil
		},
	}
}
