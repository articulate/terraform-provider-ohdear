package ohdear

import (
	"context"
	"fmt"
	"strings"

	"github.com/articulate/ohdear-sdk/ohdear"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func init() {
	schema.DescriptionKind = schema.StringMarkdown

	// add defaults on to the exported descriptions if present
	schema.SchemaDescriptionBuilder = func(s *schema.Schema) string {
		desc := s.Description
		if s.Default != nil {
			desc += fmt.Sprintf(" Defaults to `%v`.", s.Default)
		}
		return strings.TrimSpace(desc)
	}
}

func New(version string) func() *schema.Provider {
	return func() *schema.Provider {
		return Provider()
	}
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OHDEAR_TOKEN", nil),
			},
			"api_url": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("OHDEAR_API_URL", "https://ohdear.app"),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ohdear_site": resourceOhdearSite(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	apiToken := d.Get("api_token").(string)
	baseURL := d.Get("api_url").(string)

	httpClient := cleanhttp.DefaultClient()
	httpClient.Transport = logging.NewTransport("Oh Dear SDK", httpClient.Transport)
	client, err := ohdear.NewClient(baseURL, apiToken, httpClient)
	if err != nil {
		return nil, diagErrorf(err, "Unable to create Oh Dear client")
	}

	return &Config{
		apiToken: apiToken,
		baseURL:  baseURL,
		client:   client,
	}, nil
}
