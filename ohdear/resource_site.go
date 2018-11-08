package ohdear

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/articulate/ohdear-sdk/ohdear"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceOhdearSite() *schema.Resource {
	return &schema.Resource{
		Create: resourceOhdearSiteCreate,
		Read:   resourceOhdearSiteRead,
		Delete: resourceOhdearSiteDelete,
		Update: resourceOhdearSiteUpdate,
		Schema: map[string]*schema.Schema{
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
				Type:        schema.TypeSet,
				Optional:    true,
				Description: "Checks to include for side, default is all checks. Note: you cannot enable certificate checks on http URLs.",
				Elem: &schema.Schema{
					Type:         schema.TypeString,
					ValidateFunc: validation.StringInSlice(ohdear.CheckTypes, false),
				},
			},
		},
	}
}

func getSiteID(d *schema.ResourceData) (int, error) {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return id, fmt.Errorf("corrupted resource ID in terraform state, Oh Dear only supports integer IDs. Err: %v", err)
	}

	return id, err
}

func resourceOhdearSiteExists(d *schema.ResourceData, meta interface{}) (bool, error) {
	client := meta.(*Config).client
	log.Printf("[DEBUG] Calling Exists lifecycle function for site %s\n", d.Id())

	id, err := getSiteID(d)
	if err != nil {
		return false, err
	}

	if _, res, err := client.SiteService.GetSite(id); err != nil {
		if res.StatusCode == http.StatusNotFound {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceOhdearSiteCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] Calling Create lifecycle function for site")

	checks := convertInterfaceToStringArr(d.Get("checks"))
	if len(checks) == 0 {
		checks = ohdear.CheckTypes
	}

	site := &ohdear.SiteRequest{
		URL:    d.Get("url").(string),
		TeamID: d.Get("team_id").(int),
		Checks: checks,
	}

	newSite, _, err := meta.(*Config).client.SiteService.CreateSite(site)
	if err != nil {
		return fmt.Errorf("error creating site: %v", err)
	}

	d.SetId(strconv.Itoa(newSite.ID))
	return resourceOhdearSiteRead(d, meta)
}

func resourceOhdearSiteRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Calling Read lifecycle function for site %s\n", d.Id())
	id, err := getSiteID(d)
	if err != nil {
		return err
	}

	newSite, _, err := meta.(*Config).client.SiteService.GetSite(id)
	if err != nil {
		return fmt.Errorf("Failed retrieving Site: %v", err)
	}

	checks := []string{}
	for _, check := range newSite.Checks {
		if check.Enabled == true {
			checks = append(checks, check.Type)
		}
	}

	// Supporting defaulting to all enabled checks
	cfgChecks := convertInterfaceToStringArr(d.Get("checks"))
	if len(checks) != len(newSite.Checks) || len(cfgChecks) > 0 {
		d.Set("checks", checks)
	}

	d.Set("url", newSite.URL)
	d.Set("team_id", newSite.TeamID)

	return nil
}

func resourceOhdearSiteDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Calling Delete lifecycle function for site %s\n", d.Id())
	id, err := getSiteID(d)
	if err != nil {
		return err
	}

	_, err = meta.(*Config).client.SiteService.DeleteSite(id)
	return err
}

func resourceOhdearSiteUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).client
	checks := convertInterfaceToStringArrNullable(d.Get("checks"))
	id, err := getSiteID(d)
	if err != nil {
		return err
	}

	site, _, err := client.SiteService.GetSite(id)
	if err != nil {
		return err
	}

	// Sync downstream checks with config
	for _, check := range site.Checks {
		if check.Enabled {
			if !contains(checks, check.Type) {
				client.CheckService.DisableCheck(check.ID)
			}
		} else {
			if contains(checks, check.Type) {
				client.CheckService.EnableCheck(check.ID)
			}
		}
	}

	return resourceOhdearSiteRead(d, meta)
}
