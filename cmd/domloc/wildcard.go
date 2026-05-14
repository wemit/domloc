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

func wildcardCmd() *cobra.Command {
	var noHTTPS bool

	cmd := &cobra.Command{
		Use:     "wildcard <pattern> <port>",
		Short:   "Route a wildcard pattern to a local port",
		Example: `  domloc wildcard "*.foo.test" 8080`,
		Args:    cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			pattern := args[0]
			port, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("invalid port: %q", args[1])
			}

			if err := validation.Domain(pattern); err != nil {
				return err
			}
			if !validation.IsWildcard(pattern) {
				return fmt.Errorf("%q must start with *.", pattern)
			}
			if err := validation.Port(port); err != nil {
				return err
			}

			changed, err := dns.ConfigureTLDForDomain(pattern)
			if err != nil {
				return wrap("configure DNS", err)
			}
			if changed {
				ui.OK(fmt.Sprintf("DNS configured for .%s", validation.TLD(pattern)))
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
				if r.Port == port && r.Domain != pattern {
					ui.Warn(fmt.Sprintf("port %d already used by %s", port, r.Domain))
				}
			}

			reg.Add(registry.Route{
				Domain:   pattern,
				Port:     port,
				HTTPS:    useHTTPS,
				Wildcard: true,
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
			ui.OK(fmt.Sprintf("%s -> localhost:%d (%s)", pattern, port, scheme))
			return nil
		},
	}

	cmd.Flags().BoolVar(&noHTTPS, "no-https", false, "serve HTTP only; default follows domloc init setting")
	return cmd
}
