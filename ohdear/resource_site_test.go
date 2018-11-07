package ohdear

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
)

func TestAccOhdearSiteCreate(t *testing.T) {
	ri := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfigForOhdearSiteCreate(ri),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("ohdear_site.test-%d", ri), "team_id", "1853"),
					resource.TestCheckResourceAttr(fmt.Sprintf("ohdear_site.test-%d", ri), "url", fmt.Sprintf("http://www.test-%d.com", ri)),
				),
			},
		},
	})
}

func TestAccOhdearSiteLifecycle(t *testing.T) {
	ri := acctest.RandInt()
	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfigForOhdearSiteCreate(ri),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("ohdear_site.test-%d", ri), "team_id", "1853"),
					resource.TestCheckResourceAttr(fmt.Sprintf("ohdear_site.test-%d", ri), "url", fmt.Sprintf("http://www.test-%d.com", ri)),
				),
			},
			{
				Config: testConfigForOhdearSiteUpdate(ri),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(fmt.Sprintf("ohdear_site.test-%d", ri), "team_id", "1853"),
					resource.TestCheckResourceAttr(fmt.Sprintf("ohdear_site.test-%d", ri), "url", fmt.Sprintf("http://updated.test-%d.com", ri)),
				),
			},
		},
	})
}

func testConfigForOhdearSiteCreate(rInt int) string {
	return fmt.Sprintf(`
resource "ohdear_site" "test-%d" {
  team_id  = 1853
  url      = "http://www.test-%d.com"
}
`, rInt, rInt)
}

func testConfigForOhdearSiteUpdate(rInt int) string {
	return fmt.Sprintf(`
resource "ohdear_site" "test-%d" {
  team_id  = 1853
  url      = "http://updated.test-%d.com"
}`, rInt, rInt)
}
