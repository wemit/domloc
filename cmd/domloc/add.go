package main

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/wemit/domloc/internal/caddy"
	"github.com/wemit/domloc/internal/dns"
	"github.com/wemit/domloc/internal/registry"
	"github.com/wemit/domloc/internal/ui"
	"github.com/wemit/domloc/internal/validation"
)

func addCmd() *cobra.Command {
	var noHTTPS bool

	cmd := &cobra.Command{
		Use:   "add <domain> <port>",
		Short: "Route a domain to a local port",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			domain := args[0]
			port, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid port: %q", args[1])
			}

			if err := validation.Domain(domain); err != nil {
				return err
			}
			if err := validation.Port(port); err != nil {
				return err
			}

			changed, err := dns.ConfigureTLDForDomain(domain)
			if err != nil {
				return wrap("configure DNS", err)
			}
			if changed {
				ui.OK(fmt.Sprintf("DNS configured for .%s", validation.TLD(domain)))
			}

			reg, err := registry.Load()
			if err != nil {
				return wrap("load registry", err)
			}

			useHTTPS := reg.IsHTTPSDefault()
			if cmd.Flags().Changed("no-https") {
				useHTTPS = !noHTTPS
			}

			for _, r := range reg.Routes {
				if r.Port == port && r.Domain != domain {
					ui.Warn(fmt.Sprintf("port %d already used by %s", port, r.Domain))
				}
			}

			reg.Add(registry.Route{
				Domain:   domain,
				Port:     port,
				HTTPS:    useHTTPS,
				Wildcard: validation.IsWildcard(domain),
			})

			if err := reg.Save(); err != nil {
				return wrap("save registry", err)
			}
			if err := caddy.EnsureRoutes(reg.Routes); err != nil {
				return wrap("reload Caddy", err)
			}

			scheme := "https"
			if !useHTTPS {
				scheme = "http"
			}
			ui.OK(fmt.Sprintf("%s -> localhost:%d (%s)", domain, port, scheme))
			return nil
		},
	}

	cmd.Flags().BoolVar(&noHTTPS, "no-https", false, "serve HTTP only; default follows domloc init setting")
	return cmd
}
