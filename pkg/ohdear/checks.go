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
)

var AllChecks = []string{
	UptimeCheck,
	BrokenLinksCheck,
	CertificateHealthCheck,
	CertificateTransparencyCheck,
	MixedContentCheck,
	PerformanceCheck,
	DNSCheck,
}

type Check struct {
	ID      int    `json:"id"`
	Type    string `json:"type"`
	Enabled bool   `json:"enabled"`
}

func (c *Client) EnableCheck(id int) error {
	_, err := c.R().Post(fmt.Sprintf("/api/checks/%d/enable", id))
	if err != nil {
		return fmt.Errorf("could not enable check %d: %w", id, err)
	}
	return nil
}

func (c *Client) DisableCheck(id int) error {
	_, err := c.R().Post(fmt.Sprintf("/api/checks/%d/disable", id))
	if err != nil {
		return fmt.Errorf("could not disable check %d: %w", id, err)
	}
	return nil
}
