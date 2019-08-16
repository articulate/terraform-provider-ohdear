package ohdear

import (
	"github.com/articulate/ohdear-sdk/ohdear"
	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/terraform/helper/logging"
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
				DefaultFunc: schema.EnvDefaultFunc("OHDEAR_API_URL", "https://ohdear.app"),
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
	baseURL := d.Get("api_url").(string)
	httpClient := cleanhttp.DefaultClient()
	httpClient.Transport = logging.NewTransport("Oh Dear SDK", httpClient.Transport)
	client, err := ohdear.NewClient(baseURL, apiToken, httpClient)
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
