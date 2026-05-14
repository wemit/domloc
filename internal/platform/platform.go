package platform

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type OS int

const (
	MacOS OS = iota
	Linux
	Windows
	Unknown
)

func Current() OS {
	switch runtime.GOOS {
	case "darwin":
		return MacOS
	case "linux":
		return Linux
	case "windows":
		return Windows
	default:
		return Unknown
	}
}

func (o OS) String() string {
	switch o {
	case MacOS:
		return "macOS"
	case Linux:
		return "Linux"
	case Windows:
		return "Windows"
	default:
		return "Unknown"
	}
}

func IsRoot() bool {
	return os.Getuid() == 0
}

func CommandExists(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func RunSudo(args ...string) error {
	cmd := exec.Command("sudo", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func RunSudoQuiet(args ...string) error {
	cmd := exec.Command("sudo", args...)
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func RunCommand(name string, args ...string) (string, error) {
	out, err := exec.Command(name, args...).CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func ResolverDir() (string, error) {
	switch Current() {
	case MacOS:
		return "/etc/resolver", nil
	case Linux:
		return "/etc/NetworkManager/dnsmasq.d", nil
	default:
		return "", fmt.Errorf("unsupported platform: %s", Current())
	}
}

func HomebrewPrefix() string {
	if out, err := exec.Command("brew", "--prefix").Output(); err == nil {
		return strings.TrimSpace(string(out))
	}
	if runtime.GOARCH == "arm64" {
		return "/opt/homebrew"
	}
	return "/usr/local"
}
