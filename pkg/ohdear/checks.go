package ohdear

import "fmt"

const (
	UptimeCheck                  = "uptime"
	BrokenLinksCheck             = "broken_links"
	CertificateHealthCheck       = "certificate_health"
	CertificateTransparencyCheck = "certificate_transparency"
	MixedContentCheck            = "mixed_content"
	PerformanceCheck             = "performance"
	DNSCheck                     = "dns"
	LighthouseCheck              = "lighthouse"
	SitemapCheck                 = "sitemap"
	DomainCheck                  = "domain"
)

var AllChecks = []string{
	UptimeCheck,
	BrokenLinksCheck,
	CertificateHealthCheck,
	CertificateTransparencyCheck,
	MixedContentCheck,
	PerformanceCheck,
	DNSCheck,
	LighthouseCheck,
	SitemapCheck,
	DomainCheck,
}

type Check struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
}

func (c *Client) EnableCheck(id int) error {
	_, err := c.R().Post(fmt.Sprintf("/api/checks/%d/enable", id))
	return err
}

func (c *Client) DisableCheck(id int) error {
	_, err := c.R().Post(fmt.Sprintf("/api/checks/%d/disable", id))
	return err
}
