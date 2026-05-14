package main

import (
	"fmt"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/wemit/domloc/internal/registry"
)

func openCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "open <domain>",
		Short: "Open a domain in the browser",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			domain := args[0]

			reg, err := registry.Load()
			if err != nil {
				return wrap("load registry", err)
			}

			route, ok := reg.Find(domain)
			if !ok {
				return fmt.Errorf("no route for %q — run: domloc add %s <port>", domain, domain)
			}

			scheme := "https"
			if !route.HTTPS {
				scheme = "http"
			}
			url := scheme + "://" + domain

			return openBrowser(url)
		},
	}
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "linux":
		cmd = exec.Command("xdg-open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		return fmt.Errorf("unsupported platform")
	}
	return cmd.Start()
}
