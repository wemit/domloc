package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/wemit/domloc/internal/platform"
	"github.com/wemit/domloc/internal/ui"
)

func resetCmd() *cobra.Command {
	var hard bool

	cmd := &cobra.Command{
		Use:   "reset",
		Short: "Remove all domloc state and stop services",
		Long: `Stops Caddy and dnsmasq agents, removes generated configs and launchd plists.
Routes are preserved unless --hard is passed.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			ui.Header("Resetting domloc")

			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			configDir := filepath.Join(home, ".config", "domloc")

			if platform.Current() == platform.Linux {
				caddyService := "/etc/systemd/system/domloc-caddy.service"
				if _, err := os.Stat(caddyService); err == nil {
					_ = platform.RunSudoQuiet("systemctl", "disable", "--now", "domloc-caddy")
					_ = platform.RunSudoQuiet("rm", "-f", caddyService)
					_ = platform.RunSudoQuiet("systemctl", "daemon-reload")
					ui.OK("Caddy systemd service removed")
				}
			} else {
				caddyPlist := "/Library/LaunchDaemons/com.domloc.caddy.plist"
				if _, err := os.Stat(caddyPlist); err == nil {
					_ = platform.RunSudoQuiet("launchctl", "unload", caddyPlist)
					_ = platform.RunSudoQuiet("rm", "-f", caddyPlist)
					ui.OK("Caddy launchd daemon removed")
				}
			}
			_, _ = platform.RunCommand("pkill", "-f", "caddy")

			if platform.Current() == platform.Linux {
				dnsmasqAgent := filepath.Join(home, ".config", "systemd", "user", "domloc-dnsmasq.service")
				if _, err := os.Stat(dnsmasqAgent); err == nil {
					_, _ = platform.RunCommand("systemctl", "--user", "disable", "--now", "domloc-dnsmasq")
					_ = os.Remove(dnsmasqAgent)
					_, _ = platform.RunCommand("systemctl", "--user", "daemon-reload")
					ui.OK("dnsmasq systemd user service removed")
				}
			} else {
				dnsmasqAgent := filepath.Join(home, "Library", "LaunchAgents", "com.domloc.dnsmasq.plist")
				if _, err := os.Stat(dnsmasqAgent); err == nil {
					_, _ = platform.RunCommand("launchctl", "unload", dnsmasqAgent)
					_ = os.Remove(dnsmasqAgent)
					ui.OK("dnsmasq LaunchAgent removed")
				}
			}

			if platform.Current() == platform.Linux {
				resolvedDir := "/etc/systemd/resolved.conf.d"
				if entries, err := os.ReadDir(resolvedDir); err == nil {
					for _, e := range entries {
						if strings.HasPrefix(e.Name(), "domloc-") && strings.HasSuffix(e.Name(), ".conf") {
							_ = platform.RunSudoQuiet("rm", "-f", filepath.Join(resolvedDir, e.Name()))
							ui.OK(fmt.Sprintf("removed %s/%s", resolvedDir, e.Name()))
						}
					}
					_ = platform.RunSudoQuiet("systemctl", "restart", "systemd-resolved")
				}
			} else {
				resolverDir := "/etc/resolver"
				if entries, err := os.ReadDir(resolverDir); err == nil {
					for _, e := range entries {
						content, err := os.ReadFile(filepath.Join(resolverDir, e.Name()))
						if err == nil && isDomlocResolver(string(content)) {
							_ = platform.RunSudoQuiet("rm", "-f", filepath.Join(resolverDir, e.Name()))
							ui.OK(fmt.Sprintf("removed /etc/resolver/%s", e.Name()))
						}
					}
				}
			}

			toRemove := []string{
				filepath.Join(configDir, "Caddyfile"),
				filepath.Join(configDir, ".caddy-managed"),
				filepath.Join(configDir, "dnsmasq.conf"),
				filepath.Join(configDir, "dnsmasq.log"),
				filepath.Join(configDir, "caddy.log"),
			}
			for _, f := range toRemove {
				_ = os.Remove(f)
			}
			_ = os.RemoveAll(filepath.Join(configDir, "caddy-data"))
			ui.OK("config files removed")

			if hard {
				_ = os.Remove(filepath.Join(configDir, "routes.json"))
				ui.OK("routes removed")
			} else {
				ui.Warn("routes preserved — pass --hard to also remove routes.json")
			}

			ui.OK("reset complete — run `domloc init` to start fresh")
			return nil
		},
	}

	cmd.Flags().BoolVar(&hard, "hard", false, "also remove routes.json")
	return cmd
}

func isDomlocResolver(content string) bool {
	return containsStr(content, "port 5300") || containsStr(content, "nameserver 127.0.0.1")
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
