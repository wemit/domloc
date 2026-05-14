package doctor

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/wemit/domloc/internal/caddy"
	"github.com/wemit/domloc/internal/dns"
	"github.com/wemit/domloc/internal/platform"
)

type Check struct {
	Name   string
	Pass   bool
	Detail string
}

func Run() []Check {
	return []Check{
		checkOS(),
		checkBrew(),
		checkDnsmasq(),
		checkDnsmasqRunning(),
		checkCaddy(),
		checkCaddyRunning(),
		checkDNSResolution(),
	}
}

func Print(checks []Check) {
	ok := color.New(color.FgGreen).SprintFunc()
	fail := color.New(color.FgRed).SprintFunc()

	allPass := true
	for _, c := range checks {
		if c.Pass {
			fmt.Printf("  %s %s\n", ok("✓"), c.Name)
		} else {
			fmt.Printf("  %s %s", fail("✗"), c.Name)
			if c.Detail != "" {
				fmt.Printf(" — %s", c.Detail)
			}
			fmt.Println()
			allPass = false
		}
	}

	fmt.Println()
	if allPass {
		fmt.Println(ok("All checks passed."))
	} else {
		fmt.Println(fail("Some checks failed. Run `domloc init` to fix."))
	}
}

func checkOS() Check {
	p := platform.Current()
	if p == platform.MacOS || p == platform.Linux {
		return Check{Name: fmt.Sprintf("Platform (%s)", p), Pass: true}
	}
	return Check{Name: "Platform", Pass: false, Detail: fmt.Sprintf("%s not yet supported", p)}
}

func checkBrew() Check {
	if platform.CommandExists("brew") {
		return Check{Name: "Homebrew installed", Pass: true}
	}
	return Check{Name: "Homebrew installed", Pass: false, Detail: "install from https://brew.sh"}
}

func checkDnsmasq() Check {
	if dns.IsInstalled() {
		return Check{Name: "dnsmasq installed", Pass: true}
	}
	return Check{Name: "dnsmasq installed", Pass: false, Detail: "run: brew install dnsmasq"}
}

func checkDnsmasqRunning() Check {
	if dns.IsAgentRunning() {
		return Check{Name: "dnsmasq running", Pass: true}
	}
	return Check{Name: "dnsmasq running", Pass: false, Detail: "run: domloc init"}
}

func checkCaddy() Check {
	if caddy.IsInstalled() {
		return Check{Name: "Caddy installed", Pass: true}
	}
	return Check{Name: "Caddy installed", Pass: false, Detail: "run: brew install caddy"}
}

func checkCaddyRunning() Check {
	if caddy.IsRunning() {
		return Check{Name: "Caddy running", Pass: true}
	}
	return Check{Name: "Caddy running", Pass: false, Detail: "run: domloc init"}
}

func checkDNSResolution() Check {
	if err := dns.ValidateResolution("health.test"); err != nil {
		return Check{Name: "DNS resolution (.test)", Pass: false, Detail: err.Error()}
	}
	return Check{Name: "DNS resolution (.test)", Pass: true}
}
