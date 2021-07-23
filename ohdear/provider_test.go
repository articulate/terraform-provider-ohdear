package ohdear

import (
	"context"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

var testAccProvider *schema.Provider
var testAccProviderFactories = map[string]func() (*schema.Provider, error){}

func init() {
	testAccProvider = New("test")()
	testAccProviderFactories = map[string]func() (*schema.Provider, error){
		"ohdear": func() (*schema.Provider, error) {
			return New("test")(), nil
		},
	}
}

func TestProvider(t *testing.T) {
	provider := New("test")()
	if err := provider.InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ schema.Provider = *New("test")()
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
