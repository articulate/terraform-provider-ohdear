package ohdear

import (
	"github.com/articulate/ohdear-sdk/ohdear"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
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
				DefaultFunc: schema.EnvDefaultFunc("OHDEAR_API_URL", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"ohdear_site": resourceOhdearSite(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiToken := d.Get("api_token").(string)
	baseURL := "https://ohdear.app"
	client, err := ohdear.NewClient(baseURL, apiToken)
	if err != nil {
		return nil, err
	}

	config := Config{
		apiToken: apiToken,
		baseURL:  baseURL,
		client:   client,
	}
	return &config, err
}
