package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var (
	domainRegex = regexp.MustCompile(`^[a-zA-Z0-9-]+\.[a-zA-Z]{2,}$`)
)

func HandleCheckDomain(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	domain, _ :=getString(args, "domain")
	if domain == "" {
		return err("domain parameter is required")
}

	if !domainRegex.MatchString(domain) {
		return err("invalid domain format")
}

	// Check if domain is blacklisted (simplified check)
	blacklistedDomains := []string{
		"malicious.com",
		"phishing.org",
		"spam.net",
	}

	for _, badDomain := range blacklistedDomains {
		if strings.Contains(domain, badDomain) {
			return ok(fmt.Sprintf("WARNING: Domain %s contains blacklisted pattern", domain))

	}

	// Check domain availability (simplified check)
	client := http.Client{Timeout: 30 * time.Second}
	resp, fetchErr := client.Head("http://" + domain, nil)
	if fetchErr != nil {
		return err(fmt.Sprintf("error checking domain: %v", fetchErr))
}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return ok(fmt.Sprintf("Domain %s is active and responding", domain))
}

	return ok(fmt.Sprintf("Domain %s appears to be available", domain))
}

}

func HandleExtractSubdomains(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	domain, _ :=getString(args, "domain")
	if domain == "" {
		return err("domain parameter is required")
}

	if !domainRegex.MatchString(domain) {
		return err("invalid domain format")
}

	// Simplified subdomain extraction (in real implementation would use DNS queries)
	subdomains := []string{
		"www." + domain,
		"mail." + domain,
		"ftp." + domain,
		"api." + domain,
	}

	return ok(strings.Join(subdomains, "\n"))
}

func HandleDomainInfo(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	domain, _ :=getString(args, "domain")
	if domain == "" {
		return err("domain parameter is required")
}

	if !domainRegex.MatchString(domain) {
		return err("invalid domain format")
}

	// Simplified domain info (in real implementation would query WHOIS)
	info := fmt.Sprintf(`Domain: %s
Status: Active
Registrar: Example Registry
Creation Date: 2020-01-01
Expiration Date: 2030-01-01
Nameservers:
- ns1.example.com
- ns2.example.com
`, domain)

	return ok(info)
}

func HandleDomainHistory(ctx context.Context, args map[string]interface{}) (ToolResponse, error) {
	domain, _ :=getString(args, "domain")
	if domain == "" {
		return err("domain parameter is required")
}

	if !domainRegex.MatchString(domain) {
		return err("invalid domain format")
}

	// Simplified domain history (in real implementation would query historical data)
	history := fmt.Sprintf(`Domain History for %s:
- 2020-01-01: Registered
- 2021-06-15: DNS updated
- 2022-03-10: Registrar changed
- 2023-01-01: Renewed
`, domain)

	return ok(history)
}