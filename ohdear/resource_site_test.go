package ohdear

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/articulate/ohdear-sdk/ohdear"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

var checkTypesWithoutUptime = []string{
	"broken_links",
	"certificate_health",
	"mixed_content",
	"certificate_transparency",
}

var teamID string

func init() {
	teamID = os.Getenv("OHDEAR_TEAM_ID")
}

func checkImportState(s []*terraform.InstanceState) error {
	// Expect 1 site
	if len(s) != 1 {
		return fmt.Errorf("expected 1 state: %#v", s)
	}

	return nil
}

// Test Basic Creation of a Site
func TestAccOhdearSiteCreate(t *testing.T) {
	ri := acctest.RandInt()
	fqn := getTestSiteResourceFQN(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfigForOhdearSiteOneExplicitCheck(ri),
				Check: resource.ComposeTestCheckFunc(
					ensureSiteExists(fqn),
					resource.TestCheckResourceAttr(fqn, "team_id", teamID),
					resource.TestCheckResourceAttr(fqn, "url", fmt.Sprintf("https://example.org/%d", ri)),
					// Checks
					ensureChecksEnabled(fqn, ohdear.CheckTypes),
					resource.TestCheckResourceAttr(fqn, "uptime", "true"),
					resource.TestCheckResourceAttr(fqn, "broken_links", "true"),
					resource.TestCheckResourceAttr(fqn, "certificate_health", "true"),
					resource.TestCheckResourceAttr(fqn, "mixed_content", "true"),
					resource.TestCheckResourceAttr(fqn, "certificate_transparency", "true"),
				),
			},
		},
	})
}

// Test Basic Creation of a Site
func TestAccOhdearSiteCreateWithDisabledCheck(t *testing.T) {
	ri := acctest.RandInt()
	fqn := getTestSiteResourceFQN(ri)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfigForOhdearSiteUptimeDisabled(ri),
				Check: resource.ComposeTestCheckFunc(
					ensureSiteExists(fqn),
					resource.TestCheckResourceAttr(fqn, "team_id", teamID),
					resource.TestCheckResourceAttr(fqn, "url", fmt.Sprintf("https://example.org/%d", ri)),
					ensureChecksEnabled(fqn, checkTypesWithoutUptime),
					ensureChecksDisabled(fqn, []string{"uptime"}),
					resource.TestCheckResourceAttr(fqn, "uptime", "false"),
					resource.TestCheckResourceAttr(fqn, "broken_links", "true"),
					resource.TestCheckResourceAttr(fqn, "certificate_health", "true"),
					resource.TestCheckResourceAttr(fqn, "mixed_content", "true"),
					resource.TestCheckResourceAttr(fqn, "certificate_transparency", "true"),
				),
			},
		},
	})
}

func TestAccOhDearSiteCreateAddDisableThenRemoveCheckConfig(t *testing.T) {
	ri := acctest.RandInt()
	fqn := getTestSiteResourceFQN(ri)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfigForOhdearSiteNoExplicitChecks(ri),
				Check: resource.ComposeTestCheckFunc(
					ensureSiteExists(fqn),
				),
			},
			{
				Config: testConfigForOhdearSiteOneExplicitCheck(ri),
				Check: resource.ComposeTestCheckFunc(
					ensureSiteExists(fqn),
					ensureChecksEnabled(fqn, ohdear.CheckTypes),
					resource.TestCheckResourceAttr(fqn, "uptime", "true"),
					resource.TestCheckResourceAttr(fqn, "broken_links", "true"),
					resource.TestCheckResourceAttr(fqn, "certificate_health", "true"),
					resource.TestCheckResourceAttr(fqn, "mixed_content", "true"),
					resource.TestCheckResourceAttr(fqn, "certificate_transparency", "true"),
				),
			},
			{
				Config: testConfigForOhdearSiteUptimeDisabled(ri),
				Check: resource.ComposeTestCheckFunc(
					ensureSiteExists(fqn),
					ensureChecksEnabled(fqn, checkTypesWithoutUptime),
					ensureChecksDisabled(fqn, []string{"uptime"}),
					resource.TestCheckResourceAttr(fqn, "uptime", "false"),
					resource.TestCheckResourceAttr(fqn, "broken_links", "true"),
					resource.TestCheckResourceAttr(fqn, "certificate_health", "true"),
					resource.TestCheckResourceAttr(fqn, "mixed_content", "true"),
					resource.TestCheckResourceAttr(fqn, "certificate_transparency", "true"),
				),
			},
		},
	})
}

func TestAccOhdearSiteImport(t *testing.T) {
	ri := acctest.RandInt()
	fqn := getTestSiteResourceFQN(ri)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfigForOhdearSiteNoExplicitChecks(ri),
			},
			{
				ResourceName:     fqn,
				ImportState:      true,
				ImportStateCheck: checkImportState,
			},
		},
	})
}

func TestAccOhdearSiteUpdateUrl(t *testing.T) {
	ri := acctest.RandInt()
	fqn := getTestSiteResourceFQN(ri)

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfigForOhdearSiteNoExplicitChecks(ri),
				Check: resource.ComposeTestCheckFunc(
					ensureSiteExists(fqn),
					resource.TestCheckResourceAttr(fqn, "team_id", teamID),
					resource.TestCheckResourceAttr(fqn, "url", fmt.Sprintf("https://example.org/%d", ri)),
				),
			},
			{
				Config: testConfigForOhdearSiteUpdatedURL(ri),
				Check: resource.ComposeTestCheckFunc(
					ensureSiteExists(fqn),
					resource.TestCheckResourceAttr(fqn, "team_id", teamID),
					resource.TestCheckResourceAttr(fqn, "url", fmt.Sprintf("https://example.org/foo/%d", ri)),
					resource.TestCheckResourceAttr(fqn, "uptime", "true"),
					resource.TestCheckResourceAttr(fqn, "broken_links", "true"),
					resource.TestCheckResourceAttr(fqn, "certificate_health", "true"),
					resource.TestCheckResourceAttr(fqn, "mixed_content", "true"),
					resource.TestCheckResourceAttr(fqn, "certificate_transparency", "true"),
				),
			},
		},
	})
}

func TestAccOhdearSiteAddExplicitChecks(t *testing.T) {
	ri := acctest.RandInt()
	fqn := getTestSiteResourceFQN(ri)
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfigForOhdearSiteNoExplicitChecks(ri),
				Check: resource.ComposeTestCheckFunc(
					ensureSiteExists(fqn),
					resource.TestCheckResourceAttr(fqn, "uptime", "true"),
					resource.TestCheckResourceAttr(fqn, "broken_links", "true"),
					resource.TestCheckResourceAttr(fqn, "certificate_health", "true"),
					resource.TestCheckResourceAttr(fqn, "mixed_content", "true"),
					resource.TestCheckResourceAttr(fqn, "certificate_transparency", "true"),
				),
			},
			{
				Config: testConfigForOhdearSiteUptimeDisabled(ri),
				Check: resource.ComposeTestCheckFunc(
					ensureSiteExists(fqn),
					resource.TestCheckResourceAttr(fqn, "uptime", "false"),
					resource.TestCheckResourceAttr(fqn, "broken_links", "true"),
					resource.TestCheckResourceAttr(fqn, "certificate_health", "true"),
					resource.TestCheckResourceAttr(fqn, "mixed_content", "true"),
					resource.TestCheckResourceAttr(fqn, "certificate_transparency", "true"),
				),
			},
		},
	})
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("OHDEAR_TOKEN"); v == "" {
		t.Fatal("OHDEAR_TOKEN must be set for acceptance tests")
	}
	if teamID == "" {
		t.Fatal("OHDEAR_TEAM_ID must be set for acceptance tests")
	}
}

func getTestSiteResourceFQN(ri int) string {
	return fmt.Sprintf("ohdear_site.%s", getTestResourceName(ri))
}

func getTestResourceName(ri int) string {
	return fmt.Sprintf("testAcc-%d", ri)
}

func ensureChecksEnabled(name string, checksWanted []string) resource.TestCheckFunc {
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

func ensureChecksDisabled(name string, checksWanted []string) resource.TestCheckFunc {
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
			if enabled {
				return fmt.Errorf("Check %s not enabled for site %s", check, name)
			}
		}

		return nil
	}
}

func ensureSiteExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		missingErr := fmt.Errorf("resource not found: %s", name)
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return missingErr
		}
		exists, err := doesSiteExist(rs.Primary.ID)

		if !exists {
			if err != nil {
				return err
			}
			return missingErr
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

func doesSiteExist(strID string) (bool, error) {
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

func testConfigForOhdearSiteNoExplicitChecks(rInt int) string {
	name := getTestResourceName(rInt)
	return fmt.Sprintf(`
resource "ohdear_site" "%s" {
  team_id  = %s
  url      = "https://example.org/%d"
}
`, name, teamID, rInt)
}

func testConfigForOhdearSiteUpdatedURL(rInt int) string {
	name := getTestResourceName(rInt)
	return fmt.Sprintf(`
resource "ohdear_site" "%s" {
  team_id  = %s
  url      = "https://example.org/foo/%d"
}`, name, teamID, rInt)
}

func testConfigForOhdearSiteOneExplicitCheck(rInt int) string {
	name := getTestResourceName(rInt)
	return fmt.Sprintf(`
resource "ohdear_site" "%s" {
  team_id  = %s
  url      = "https://example.org/%d"

  uptime = true
}`, name, teamID, rInt)
}

func testConfigForOhdearSiteUptimeDisabled(rInt int) string {
	name := getTestResourceName(rInt)
	return fmt.Sprintf(`
resource "ohdear_site" "%s" {
  team_id  = %s
  url      = "https://example.org/%d"

  uptime = false
}`, name, teamID, rInt)
}
