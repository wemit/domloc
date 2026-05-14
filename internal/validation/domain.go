package validation

import (
	"fmt"
	"regexp"
	"strings"
)

var domainRe = regexp.MustCompile(`^(\*\.)?([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)

func Domain(domain string) error {
	if domain == "" {
		return fmt.Errorf("domain cannot be empty")
	}
	if !domainRe.MatchString(domain) {
		return fmt.Errorf("invalid domain: %q", domain)
	}
	return nil
}

func Port(port int) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("port must be 1-65535, got %d", port)
	}
	return nil
}

func IsWildcard(domain string) bool {
	return strings.HasPrefix(domain, "*.")
}

func TLD(domain string) string {
	parts := strings.Split(strings.TrimPrefix(domain, "*."), ".")
	return parts[len(parts)-1]
}
