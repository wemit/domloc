package main

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/wemit/domloc/internal/caddy"
	"github.com/wemit/domloc/internal/registry"
)

func listCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all domain routes",
		RunE: func(cmd *cobra.Command, args []string) error {
			reg, err := registry.Load()
			if err != nil {
				return wrap("load registry", err)
			}

			if len(reg.Routes) == 0 {
				fmt.Println("No routes. Add one with: domloc add <domain> <port>")
				return nil
			}

			running := caddy.IsRunning()
			bold := color.New(color.Bold).SprintFunc()
			dim := color.New(color.Faint).SprintFunc()

			fmt.Printf("  %s %s %s %s %s\n",
				pad(bold("DOMAIN"), "DOMAIN", 30),
				pad(bold("PORT"), "PORT", 8),
				pad(bold("HTTPS"), "HTTPS", 7),
				pad(bold("WILDCARD"), "WILDCARD", 10),
				bold("PROXY"))

			for _, r := range reg.Routes {
				httpsStr, httpsRaw := dim("no"), "no"
				if r.HTTPS {
					httpsStr, httpsRaw = color.GreenString("yes"), "yes"
				}
				wildcardStr, wildcardRaw := dim("no"), "no"
				if r.Wildcard {
					wildcardStr, wildcardRaw = color.CyanString("yes"), "yes"
				}
				proxyStr := dim("stopped")
				if running {
					proxyStr = color.GreenString("running")
				}

				fmt.Printf("  %-30s %-8d %s %s %s\n",
					r.Domain,
					r.Port,
					pad(httpsStr, httpsRaw, 7),
					pad(wildcardStr, wildcardRaw, 10),
					proxyStr,
				)
			}

			return nil
		},
	}
}

// pad appends spaces so the visible width of a colored string matches width.
func pad(colored, raw string, width int) string {
	spaces := width - len(raw)
	if spaces < 0 {
		spaces = 0
	}
	return colored + strings.Repeat(" ", spaces)
}
