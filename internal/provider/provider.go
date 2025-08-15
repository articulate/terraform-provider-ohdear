package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/articulate/terraform-provider-ohdear/internal/runtime"
	"github.com/articulate/terraform-provider-ohdear/pkg/ohdear"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown

	// add defaults on to the exported descriptions if present
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		if s.Deprecated != "" {
			desc += " __Deprecated__: " + s.Deprecated
		}
		return strings.TrimSpace(desc)
	}
}

func New() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Oh Dear API token. If not set, uses `OHDEAR_TOKEN` env var.",
				DefaultFunc: schema.EnvDefaultFunc("OHDEAR_TOKEN", nil),
			},
			"api_url": {
				Type:         schema.TypeString,
				Optional:     true,
				Description:  "Oh Dear API URL. If not set, uses `OHDEAR_API_URL` env var. Defaults to `https://ohdear.app`.",
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
				DefaultFunc:  schema.EnvDefaultFunc("OHDEAR_API_URL", "https://ohdear.app"),
			},
			"team_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "The default team ID to use for sites. If not set, uses `OHDEAR_TEAM_ID` env var.",
				DefaultFunc: schema.EnvDefaultFunc("OHDEAR_TEAM_ID", 0),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ohdear_site":    resourceOhdearSite(),
			"ohdear_monitor": resourceOhdearMonitor(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(_ context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	ua := fmt.Sprintf(
		"terraform-provider-ohdear/%s (https://github.com/articulate/terraform-provider-ohdear)",
		runtime.Version,
	)
	client := ohdear.NewClient(d.Get("api_url").(string), d.Get("api_token").(string))
	client.SetUserAgent(ua)

	return &Config{
		client: client,
		teamID: d.Get("team_id").(int),
	}, nil
}
