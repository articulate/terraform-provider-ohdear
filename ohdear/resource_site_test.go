package ohdear

import (
	"fmt"
	"github.com/hashicorp/terraform/terraform"
	"net/http"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOhdearSiteCreate(t *testing.T) {
	ri := acctest.RandInt()
	fqn := getTestSiteResourceFQN(ri)
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfigForOhdearSiteCreate(ri),
				Check: resource.ComposeTestCheckFunc(
					ensureSiteExists(fqn),
					resource.TestCheckResourceAttr(fqn, "team_id", "11"),
					resource.TestCheckResourceAttr(fqn, "url", "https://www.google.com"),
				),
			},
		},
	})
}

func TestAccOhdearSiteLifecycle(t *testing.T) {
	ri := acctest.RandInt()
	fqn := getTestSiteResourceFQN(ri)
	resource.Test(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: ensureSiteDestroyed,
		Steps: []resource.TestStep{
			{
				Config: testConfigForOhdearSiteCreate(ri),
				Check: resource.ComposeTestCheckFunc(
					ensureSiteExists(fqn),
					resource.TestCheckResourceAttr(fqn, "team_id", "11"),
					resource.TestCheckResourceAttr(fqn, "url", "https://www.google.com"),
				),
			},
			{
				Config: testConfigForOhdearSiteUpdate(ri),
				Check: resource.ComposeTestCheckFunc(
					ensureSiteExists(fqn),
					resource.TestCheckResourceAttr(fqn, "team_id", "11"),
					resource.TestCheckResourceAttr(fqn, "url", "https://www.bing.com"),
					resource.TestCheckResourceAttr(fqn, "checks.#", "1"),
				),
			},
		},
	})
}

func getTestSiteResourceFQN(ri int) string {
	return fmt.Sprintf("ohdear_site.%s", getTestResourceName(ri))
}

func getTestResourceName(ri int) string {
	return fmt.Sprintf("testAcc-%d", ri)
}

func ensureSiteDestroyed(s *terraform.State) error {
	for _, r := range s.RootModule().Resources {
		exists, err := doesSiteExist(r.Primary.ID)
		if exists {
			if err != nil {
			}

			return fmt.Errorf("Test site still exists, beware of the danglers")
		}
		return err
	}

	return nil
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

func testConfigForOhdearSiteCreate(rInt int) string {
	name := getTestResourceName(rInt)
	return fmt.Sprintf(`
resource "ohdear_site" "%s" {
  team_id  = 11
  url      = "https://www.google.com"
}
`, name)
}

func testConfigForOhdearSiteUpdate(rInt int) string {
	name := getTestResourceName(rInt)
	return fmt.Sprintf(`
resource "ohdear_site" "%s" {
  team_id  = 11
  url      = "https://www.bing.com"
  checks   = [
	  "uptime"
  ]
}`, name)
}
