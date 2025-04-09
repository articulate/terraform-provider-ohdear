package provider

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var (
	testAccProvider          *schema.Provider
	testAccProviderFactories = map[string]func() (*schema.Provider, error){}
)

func init() {
	testAccProvider = New()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"ohdear": func() (*schema.Provider, error) {
			return New(), nil
		},
	}
}

func TestProvider(t *testing.T) {
	provider := New()
	if err := provider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(_ *testing.T) {
	var _ schema.Provider = *New() //nolint:staticcheck
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("OHDEAR_TOKEN"); v == "" {
		t.Fatal("OHDEAR_TOKEN must be set for acceptance tests")
	}
	if teamID == "" {
		t.Fatal("OHDEAR_TEAM_ID must be set for acceptance tests")
	}

	diags := testAccProvider.Configure(context.TODO(), terraform.NewResourceConfigRaw(nil))
	if diags.HasError() {
		t.Fatal(diags[0].Summary)
	}
}
