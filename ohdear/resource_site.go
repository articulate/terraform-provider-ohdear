package ohdear

import (
	"github.com/hashicorp/terraform/helper/schema"
	//"github.com/hashicorp/terraform/helper/validation"
)

func resourceSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceSiteCreate,
		Read:   resourceSiteRead,
		Update: resourceSiteUpdate,
		Delete: resourceSiteDelete,
		Schema: map[string]*schema.Schema{
			"id": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ID of the site - assigned by OhDear",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "URL of the site to be checked",
			},
		},
	}
}

func resourceSiteExists(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSiteCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSiteUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSiteRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSiteDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
