package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/wemit/domloc/internal/caddy"
	"github.com/wemit/domloc/internal/dns"
	"github.com/wemit/domloc/internal/platform"
	"github.com/wemit/domloc/internal/registry"
	"github.com/wemit/domloc/internal/services"
	"github.com/wemit/domloc/internal/ui"
)

func initCmd() *cobra.Command {
	var noHTTPS bool

	cmd := &cobra.Command{
		Use:   "init",
		Short: "Install dependencies and start Caddy",
		RunE: func(cmd *cobra.Command, args []string) error {
			ui.Header("Initializing domloc")

			p := platform.Current()
			if p != platform.MacOS && p != platform.Linux {
				ui.Fatalf("only macOS and Linux are supported")
			}

			if !dns.IsInstalled() {
				ui.Warn("dnsmasq not found, installing...")
				if err := dns.Install(); err != nil {
					return wrap("install dnsmasq", err)
				}
			}
			ui.OK("dnsmasq installed")

			if !dns.IsAgentRunning() {
				if err := dns.InstallAgent(); err != nil {
					return wrap("start dnsmasq agent", err)
				}
			}
			ui.OK("dnsmasq running")

			if !caddy.IsInstalled() {
				ui.Warn("Caddy not found, installing...")
				if err := caddy.Install(); err != nil {
					return wrap("install caddy", err)
				}
			}
			ui.OK("Caddy installed")

			reg, err := registry.Load()
			if err != nil {
				return wrap("load registry", err)
			}

			reg.HTTPSDefault = registry.BoolPtr(!noHTTPS)
			if err := reg.Save(); err != nil {
				return wrap("save registry", err)
			}

			if caddy.IsAdminAPIRunning() && !services.IsInstalled() {
				ui.Warn("existing Caddy detected — injecting routes via admin API (your config untouched)")
			}

			if err := caddy.EnsureRoutes(reg.Routes); err != nil {
				printCaddyLog()
				return wrap("start Caddy", err)
			}
			ui.OK("Caddy running")

			if noHTTPS {
				ui.Warn("HTTPS disabled — routes will serve HTTP only")
			} else {
				if err := caddy.TrustLocalCA(); err != nil {
					ui.Warn("could not auto-trust local CA — run manually: sudo caddy trust")
				} else {
					ui.OK("local HTTPS CA trusted")
				}
			}

			ui.OK("domloc ready — run `domloc add <domain> <port>` to get started")
			return nil
		},
	}

	cmd.Flags().BoolVar(&noHTTPS, "no-https", false, "skip CA trust step (no sudo for keychain, HTTP-only routes by default)")
	return cmd
}

func printCaddyLog() {
	home, _ := os.UserHomeDir()
	logPath := filepath.Join(home, ".config", "domloc", "caddy.log")
	data, err := os.ReadFile(logPath)
	if err != nil || len(data) == 0 {
		return
	}
	fmt.Fprintf(os.Stderr, "\n--- caddy.log ---\n%s\n--- end ---\n\n", data)
}

func wrap(op string, err error) error {
	if err == nil {
		return nil
	}
	return &opError{op: op, err: err}
}

type opError struct {
	op  string
	err error
}

func (e *opError) Error() string {
	return e.op + ": " + e.err.Error()
}

func (e *opError) Unwrap() error {
	return e.err
}
