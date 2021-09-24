package ohdear

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/articulate/ohdear-sdk/ohdear"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var teamID string

func init() {
	teamID = os.Getenv("OHDEAR_TEAM_ID")
}

func TestAccOhdearSite(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-test")
	url := fmt.Sprintf("https://example.com/%s", name)
	resourceName := fmt.Sprintf("ohdear_site.%s", name)
	updatedURL := fmt.Sprintf("%s/new", url)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		IDRefreshName:     resourceName,
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOhdearSiteConfigBasic(name, url),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEnsureSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "team_id", teamID),
					resource.TestCheckResourceAttr(resourceName, "url", url),
					testAccEnsureChecksEnabled(resourceName, ohdear.CheckTypes),
					testAccEnsureChecksEnabled(resourceName, []string{"performance"}),
					resource.TestCheckResourceAttr(resourceName, "checks.0.uptime", "true"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.broken_links", "true"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.certificate_health", "true"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.certificate_transparency", "true"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.mixed_content", "true"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.performance", "true"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOhdearSiteConfigBasic(name, updatedURL),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEnsureSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "team_id", teamID),
					resource.TestCheckResourceAttr(resourceName, "url", updatedURL),
				),
			},
		},
	})
}

func TestAccOhdearSite_EnableDisableChecks(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-test")
	url := fmt.Sprintf("https://example.com/%s", name)
	resourceName := fmt.Sprintf("ohdear_site.%s", name)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		IDRefreshName:     resourceName,
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOhdearSiteConfigChecks(name, url, map[string]bool{"uptime": true, "broken_links": true}),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEnsureSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "team_id", teamID),
					resource.TestCheckResourceAttr(resourceName, "url", url),
					resource.TestCheckResourceAttr(resourceName, "checks.0.uptime", "true"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.broken_links", "true"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.certificate_health", "false"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.mixed_content", "false"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.certificate_transparency", "false"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.performance", "false"),
					testAccEnsureChecksEnabled(resourceName, []string{"uptime", "broken_links"}),
					testAccEnsureChecksDisabled(resourceName, []string{"mixed_content", "performance"}),
				),
			},
			{
				Config: testAccOhdearSiteConfigChecks(name, url, map[string]bool{"uptime": true}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "url", url),
					resource.TestCheckResourceAttr(resourceName, "checks.0.uptime", "true"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.broken_links", "false"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.performance", "false"),
					testAccEnsureChecksEnabled(resourceName, []string{"uptime"}),
					testAccEnsureChecksDisabled(resourceName, []string{"broken_links", "performance"}),
				),
			},
			{
				Config: testAccOhdearSiteConfigChecks(name, url, map[string]bool{"uptime": false}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "url", url),
					resource.TestCheckResourceAttr(resourceName, "checks.0.uptime", "false"),
					testAccEnsureChecksDisabled(resourceName, []string{"uptime", "broken_links"}),
				),
			},
		},
	})
}

func TestAccOhdearSite_TeamID(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-test")
	url := fmt.Sprintf("https://example.com/%s", name)
	resourceName := fmt.Sprintf("ohdear_site.%s", name)

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		IDRefreshName:     resourceName,
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckSiteDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOhdearSiteConfigBasic(name, url),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEnsureSiteExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "team_id", teamID),
					resource.TestCheckResourceAttr(resourceName, "url", url),
				),
			},
			{
				Config:             testAccOhdearSiteConfigTeamID(name, url, "1"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "team_id", "1"),
				),
			},
		},
	})
}

// Checks

func doesSiteExists(strID string) (bool, error) {
	client := testAccProvider.Meta().(*Config).client
	id, _ := strconv.Atoi(strID)
	if _, res, err := client.SiteService.GetSite(id); err != nil {
		if res.StatusCode == http.StatusNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func testAccCheckSiteDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ohdear_site" {
			continue
		}

		// give the API time to update
		time.Sleep(5 * time.Second)

		exists, err := doesSiteExists(rs.Primary.ID)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("site still exists in Oh Dear: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccEnsureSiteExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}

		exists, err := doesSiteExists(rs.Primary.ID)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("resource not found: %s", name)
		}

		return nil
	}
}

func testAccEnsureChecksEnabled(name string, checksWanted []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Config).client

		missingErr := fmt.Errorf("resource not found: %s", name)
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return missingErr
		}
		siteID, _ := strconv.Atoi(rs.Primary.ID)
		site, _, _ := client.SiteService.GetSite(siteID)

		for _, check := range checksWanted {
			enabled := isCheckEnabled(site, check)
			if !enabled {
				return fmt.Errorf("Check %s not enabled for site %s", check, name)
			}
		}

		return nil
	}
}

func testAccEnsureChecksDisabled(name string, checksWanted []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*Config).client

		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}

		siteID, _ := strconv.Atoi(rs.Primary.ID)
		site, _, err := client.SiteService.GetSite(siteID)
		if err != nil {
			return err
		}

		for _, check := range checksWanted {
			if isCheckEnabled(site, check) {
				return fmt.Errorf("check %s not enabled for site %s", check, name)
			}
		}

		return nil
	}
}

// isCheckEnabled checks the site retrieved from OhDear to see whether the
// specified check is present and enabled
func isCheckEnabled(site *ohdear.Site, checkName string) bool {
	for _, aCheck := range site.Checks {
		if aCheck.Type == checkName && aCheck.Enabled == true {
			return true
		}
	}

	return false
}

// Configs

func testAccOhdearSiteConfigBasic(name, url string) string {
	return fmt.Sprintf(`
resource "ohdear_site" "%s" {
  url = "%s"
}
`, name, url)
}

func testAccOhdearSiteConfigTeamID(name, url, team string) string {
	return fmt.Sprintf(`
resource "ohdear_site" "%s" {
  team_id = %s
  url     = "%s"
}
`, name, team, url)
}

func testAccOhdearSiteConfigChecks(name, url string, checks map[string]bool) string {
	block := []string{}
	for check, enabled := range checks {
		block = append(block, fmt.Sprintf("%s = %t", check, enabled))
	}

	return fmt.Sprintf(`
resource "ohdear_site" "%s" {
  url = "%s"

  checks {
    %s
  }
}
`, name, url, strings.Join(block, "\n    "))
}
