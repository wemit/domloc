package services

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/wemit/domloc/internal/platform"
)

const (
	label    = "com.domloc.caddy"
	plistDst = "/Library/LaunchDaemons/com.domloc.caddy.plist"
	adminAPI = "http://localhost:2019"
)

const plistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.domloc.caddy</string>
	<key>ProgramArguments</key>
	<array>
		<string>{{.CaddyBin}}</string>
		<string>run</string>
		<string>--config</string>
		<string>{{.Caddyfile}}</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<true/>
	<key>StandardOutPath</key>
	<string>{{.LogFile}}</string>
	<key>StandardErrorPath</key>
	<string>{{.LogFile}}</string>
</dict>
</plist>
`

type plistData struct {
	CaddyBin  string
	Caddyfile string
	LogFile   string
}

func Install(caddyfilePath string) error {
	if platform.Current() == platform.Linux {
		return installLinux(caddyfilePath)
	}
	caddyBin, err := findCaddy()
	if err != nil {
		return err
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	logFile := filepath.Join(home, ".config", "domloc", "caddy.log")

	tmpl, err := template.New("plist").Parse(plistTemplate)
	if err != nil {
		return err
	}

	tmp, err := os.CreateTemp("", "domloc-caddy-*.plist")
	if err != nil {
		return err
	}
	defer os.Remove(tmp.Name())

	if err := tmpl.Execute(tmp, plistData{
		CaddyBin:  caddyBin,
		Caddyfile: caddyfilePath,
		LogFile:   logFile,
	}); err != nil {
		return err
	}
	tmp.Close()

	if err := platform.RunSudo("cp", tmp.Name(), plistDst); err != nil {
		return fmt.Errorf("copy plist: %w", err)
	}
	if err := platform.RunSudo("chown", "root:wheel", plistDst); err != nil {
		return fmt.Errorf("chown plist: %w", err)
	}

	_ = platform.RunSudoQuiet("launchctl", "unload", plistDst)
	if err := platform.RunSudo("launchctl", "load", "-w", plistDst); err != nil {
		return fmt.Errorf("launchctl load: %w", err)
	}

	return waitForAdminAPI(15 * time.Second)
}

func Reload(caddyfilePath string) error {
	if platform.Current() == platform.Linux {
		return reloadLinux(caddyfilePath)
	}
	_, err := platform.RunCommand("caddy", "reload", "--config", caddyfilePath)
	return err
}

func Uninstall() error {
	if platform.Current() == platform.Linux {
		return uninstallLinux()
	}
	_ = platform.RunSudoQuiet("launchctl", "unload", plistDst)
	_ = platform.RunSudoQuiet("rm", "-f", plistDst)
	return nil
}

func IsInstalled() bool {
	if platform.Current() == platform.Linux {
		return isInstalledLinux()
	}
	_, err := os.Stat(plistDst)
	return err == nil
}

func IsRunning() bool {
	if platform.Current() == platform.Linux {
		return isRunningLinux()
	}
	out, err := platform.RunCommand("launchctl", "list", label)
	if err != nil {
		return false
	}
	return strings.Contains(out, label)
}

func waitForAdminAPI(timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 300 * time.Millisecond}
	for time.Now().Before(deadline) {
		resp, err := client.Get(adminAPI + "/config/")
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == 200 {
				return nil
			}
		}
		time.Sleep(250 * time.Millisecond)
	}
	return fmt.Errorf("caddy admin API did not become ready within %s", timeout)
}

func findCaddy() (string, error) {
	out, err := platform.RunCommand("which", "caddy")
	if err != nil || strings.TrimSpace(out) == "" {
		return "", fmt.Errorf("caddy not found in PATH")
	}
	return strings.TrimSpace(out), nil
}
