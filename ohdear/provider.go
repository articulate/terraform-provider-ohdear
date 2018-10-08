package ohdear

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_token": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				EnvDefaultFunc("OHDEAR_TOKEN", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"site": resourceSite(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		Token: d.Get("api_token").(string),
	}
}
