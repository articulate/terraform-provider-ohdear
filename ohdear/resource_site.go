package ohdear

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/articulate/ohdear-sdk/ohdear"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOhdearSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceOhdearSiteCreate,
		Read:   resourceOhdearSiteRead,
		Delete: resourceOhdearSiteDelete,
		Schema: map[string]*schema.Schema{
			"site_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Primary Key of the site",
			},
			"url": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "URL of the site to be checked",
			},
			"team_id": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "ID of the team for this site",
			},
			"checks": &schema.Schema{
				Type:     schema.TypeSet,
				ForceNew: true,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uptime": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"broken_links": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"certificate_health": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"mixed_content": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
						"certificate_transparency": &schema.Schema{
							Type:     schema.TypeBool,
							Optional: true,
							Default:  true,
						},
					},
				},
			},
		},
	}
}

func resourceOhdearSiteExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*Config).client
	log.Printf("[DEBUG] Calling Exists lifecycle function for site %v\n", d.Id)
	if _, _, err := client.SiteService.GetSite(d.Get("site_id").(int)); err != nil {
		return false, err
	}

	return true, nil
}

func resourceOhdearSiteCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Calling Create lifecycle function for site %v\n", d.Id)
	site := &ohdear.Site{
		Url:    d.Get("url").(string),
		TeamID: d.Get("team_id").(int),
	}

	newSite, _, err := meta.(*Config).client.SiteService.CreateSite(site)

	if err != nil {
		return fmt.Errorf("error creating site: %s", err.Error())
	}

	d.Set("site_id", newSite.ID)
	d.SetId(d.Get("url").(string))
	return resourceOhdearSiteRead(d, meta)
}

func resourceOhdearSiteRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Calling Read lifecycle function for site %v\n", d.Id)
	id := d.Get("site_id").(int)
	newSite, resp, err := meta.(*Config).client.SiteService.GetSite(id)

	defer resp.Body.Close()

	htmlData, err := ioutil.ReadAll(resp.Body)
	log.Printf("[DEBUG] Response from Ohdear API: %s", string(htmlData))
	if err != nil {
		return fmt.Errorf("Error reading Site: %s", err.Error())
	}

	d.Set("url", newSite.Url)
	d.Set("team_id", newSite.TeamID)

	return nil
}

func resourceOhdearSiteDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Calling Delete lifecycle function for site %v\n", d.Id)
	id := d.Get("site_id").(int)

	_, err := meta.(*Config).client.SiteService.DeleteSite(id)
	return err
}
