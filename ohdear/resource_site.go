package ohdear

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteCreate,
		Read:   resourceSiteRead,
		Update: resourceSiteUpdate,
		Delete: resourceSiteDelete,
		Schema: map[string]*schema.Schema{
			"id": &Schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the site - assigned by OhDear",
			},
			"url": &Schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL of the site to be checked",
			},
		},
	}
}
