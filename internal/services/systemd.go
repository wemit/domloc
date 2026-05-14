package services

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/wemit/domloc/internal/platform"
)

const (
	linuxServiceName = "domloc-caddy"
	linuxUnitPath    = "/etc/systemd/system/domloc-caddy.service"
)

const systemdUnitTemplate = `[Unit]
Description=Caddy reverse proxy (domloc)
After=network.target

[Service]
ExecStart={{.CaddyBin}} run --config {{.Caddyfile}}
ExecReload=/bin/kill -USR1 $MAINPID
Restart=always
RestartSec=2
StandardOutput=append:{{.LogFile}}
StandardError=append:{{.LogFile}}

[Install]
WantedBy=multi-user.target
`

type systemdUnitData struct {
	CaddyBin  string
	Caddyfile string
	LogFile   string
}

func installLinux(caddyfilePath string) error {
	caddyBin, err := findCaddy()
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	logFile := filepath.Join(home, ".config", "domloc", "caddy.log")

	tmpl, err := template.New("systemd").Parse(systemdUnitTemplate)
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp("", "domloc-caddy-*.service")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if err := tmpl.Execute(tmp, systemdUnitData{
		CaddyBin:  caddyBin,
		Caddyfile: caddyfilePath,
		LogFile:   logFile,
	}); err != nil {
		return err
	}
	tmp.Close()

	if err := platform.RunSudo("cp", tmp.Name(), linuxUnitPath); err != nil {
		return fmt.Errorf("copy unit file: %w", err)
	}

	if err := platform.RunSudo("systemctl", "daemon-reload"); err != nil {
		return fmt.Errorf("daemon-reload: %w", err)
	}

	if err := platform.RunSudo("systemctl", "enable", "--now", linuxServiceName); err != nil {
		return fmt.Errorf("systemctl enable: %w", err)
	}

	return waitForAdminAPI(15 * time.Second)
}

func reloadLinux(caddyfilePath string) error {
	_, err := platform.RunCommand("caddy", "reload", "--config", caddyfilePath)
	return err
}

func uninstallLinux() error {
	_ = platform.RunSudoQuiet("systemctl", "disable", "--now", linuxServiceName)
	_ = platform.RunSudoQuiet("rm", "-f", linuxUnitPath)
	_ = platform.RunSudoQuiet("systemctl", "daemon-reload")
	return nil
}

func isInstalledLinux() bool {
	_, err := os.Stat(linuxUnitPath)
	return err == nil
}

func isRunningLinux() bool {
	out, err := platform.RunCommand("systemctl", "is-active", linuxServiceName)
	return err == nil && strings.TrimSpace(out) == "active"
}
