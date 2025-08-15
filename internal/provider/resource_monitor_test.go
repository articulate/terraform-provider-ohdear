package provider

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/articulate/terraform-provider-ohdear/pkg/ohdear"
)

func TestAccOhdearMonitor(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-test")
	url := "https://example.com/" + name
	resourceName := "ohdear_monitor." + name
	updatedURL := url + "/new"

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		IDRefreshName:     resourceName,
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOhdearMonitorConfigBasic(name, url),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEnsureMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "team_id", teamID),
					resource.TestCheckResourceAttr(resourceName, "url", url),
					testAccEnsureChecksEnabled(resourceName, []string{
						"uptime", "broken_links", "certificate_health",
						"mixed_content", "performance",
					}),
					testAccEnsureChecksDisabled(resourceName, []string{"dns"}),
					resource.TestCheckResourceAttr(resourceName, "checks.0.uptime", "true"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.broken_links", "true"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.certificate_health", "true"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.mixed_content", "true"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.performance", "true"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.dns", "false"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccOhdearMonitorConfigBasic(name, updatedURL),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEnsureMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "team_id", teamID),
					resource.TestCheckResourceAttr(resourceName, "url", updatedURL),
				),
			},
		},
	})
}

func TestAccOhdearMonitor_EnableDisableChecks(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-test")
	url := "https://example.com/" + name
	resourceName := "ohdear_monitor." + name

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		IDRefreshName:     resourceName,
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOhdearMonitorConfigChecks(name, url, map[string]bool{"uptime": true, "broken_links": true}),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEnsureMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "team_id", teamID),
					resource.TestCheckResourceAttr(resourceName, "url", url),
					resource.TestCheckResourceAttr(resourceName, "checks.0.uptime", "true"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.broken_links", "true"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.certificate_health", "false"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.mixed_content", "false"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.performance", "true"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.dns", "false"),
					testAccEnsureChecksEnabled(resourceName, []string{"uptime", "broken_links"}),
					testAccEnsureChecksDisabled(resourceName, []string{"mixed_content"}),
				),
			},
			{
				Config: testAccOhdearMonitorConfigChecks(name, url, map[string]bool{"uptime": true}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "url", url),
					resource.TestCheckResourceAttr(resourceName, "checks.0.uptime", "true"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.broken_links", "false"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.performance", "true"),
					testAccEnsureChecksEnabled(resourceName, []string{"uptime"}),
					testAccEnsureChecksDisabled(resourceName, []string{"broken_links"}),
				),
			},
			{
				Config: testAccOhdearMonitorConfigChecks(name, url, map[string]bool{"uptime": false}),
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "url", url),
					resource.TestCheckResourceAttr(resourceName, "checks.0.uptime", "false"),
					testAccEnsureChecksDisabled(resourceName, []string{"uptime", "broken_links"}),
				),
			},
		},
	})
}

func TestAccOhdearMonitor_TeamID(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-test")
	url := "https://example.com/" + name
	resourceName := "ohdear_monitor." + name

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		IDRefreshName:     resourceName,
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccOhdearMonitorConfigBasic(name, url),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccEnsureMonitorExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "team_id", teamID),
					resource.TestCheckResourceAttr(resourceName, "url", url),
				),
			},
			{
				Config:             testAccOhdearMonitorConfigTeamID(name, url, "1"),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "team_id", "1"),
				),
			},
		},
	})
}

func TestAccOhdearMonitor_HTTPDefaults(t *testing.T) {
	name := acctest.RandomWithPrefix("tf-acc-test")
	url := "http://example.com/" + name
	resourceName := "ohdear_monitor." + name

	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		IDRefreshName:     resourceName,
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckMonitorDestroy,
		Steps: []resource.TestStep{
			{
				Config:             testAccOhdearMonitorConfigBasic(name, url),
				PlanOnly:           true,
				ExpectNonEmptyPlan: true,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "checks.0.uptime", "true"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.broken_links", "true"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.certificate_health", "false"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.mixed_content", "false"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.performance", "true"),
					resource.TestCheckResourceAttr(resourceName, "checks.0.dns", "false"),
				),
			},
		},
	})
}

// Checks

func doesMonitorExist(strID string) (bool, error) {
	client := testAccProvider.Meta().(*Config).client
	id, _ := strconv.Atoi(strID)
	if _, err := client.GetMonitor(id); err != nil {
		var e *ohdear.Error
		if errors.As(err, &e) && e.Response.StatusCode() == 404 {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func testAccCheckMonitorDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ohdear_monitor" {
			continue
		}

		// give the API time to update
		time.Sleep(5 * time.Second)

		exists, err := doesMonitorExist(rs.Primary.ID)
		if err != nil {
			return err
		}

		if exists {
			return fmt.Errorf("monitor still exists in Oh Dear: %s", rs.Primary.ID)
		}
	}
	return nil
}

func testAccEnsureMonitorExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("resource not found: %s", name)
		}

		exists, err := doesMonitorExist(rs.Primary.ID)
		if err != nil {
			return err
		}

		if !exists {
			return fmt.Errorf("resource not found: %s", name)
		}

		return nil
	}
}

// Configs

func testAccOhdearMonitorConfigBasic(name, url string) string {
	return fmt.Sprintf(`
resource "ohdear_monitor" "%s" {
  url = "%s"
}
`, name, url)
}

func testAccOhdearMonitorConfigTeamID(name, url, team string) string {
	return fmt.Sprintf(`
resource "ohdear_monitor" "%s" {
  team_id = %s
  url     = "%s"
}
`, name, team, url)
}

func testAccOhdearMonitorConfigChecks(name, url string, checks map[string]bool) string {
	block := []string{}
	for check, enabled := range checks {
		block = append(block, fmt.Sprintf("%s = %t", check, enabled))
	}

	return fmt.Sprintf(`
resource "ohdear_monitor" "%s" {
  url = "%s"

  checks {
    %s
  }
}
`, name, url, strings.Join(block, "\n    "))
}
