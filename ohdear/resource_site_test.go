package ohdear

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/acctest"
	"github.com/hashicorp/terraform/helper/resource"
	"testing"
)

func TestAccOhdearSiteCreate(t *testing.T) {
	ri := acctest.RandInt()

	resource.Test(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testConfigForOhdearSiteCreate(ri),
			},
		},
	})
}

func testConfigForOhdearSiteCreate(rInt int) string {
	return fmt.Sprintf(`
resource "site" "test-%d" {
  team_id  = "1775"
  url      = "http://www.test-%d.com"
}
`, rInt, rInt)
}
