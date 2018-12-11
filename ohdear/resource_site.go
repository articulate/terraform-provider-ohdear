package ohdear

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/articulate/ohdear-sdk/ohdear"
	"github.com/hashicorp/terraform/helper/schema"
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
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Checks to include or exclude for site. Note: you cannot enable certificate checks on http URLs.",
				Elem:        schema.TypeBool,
				Computed:    true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
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

// determineChecksWanted returns the types of checks which are specified
// explicitly as enabled in our config OR which are implicitly enabled
// by virtue of their exclusion from config.
func determineChecksWanted(d *schema.ResourceData) ([]string, error) {
	// We want all the checks by default...
	checksWanted := make([]string, len(ohdear.CheckTypes))
	copy(checksWanted, ohdear.CheckTypes)

	checksConfig := d.Get("checks").(map[string]interface{})
	checksInConfig := getKeysAsSlice(d.Get("checks").(map[string]interface{}))
	numChecks := len(checksInConfig)

	// If we specified checks in the config...
	if numChecks > 0 {
		// Ensure all checks specified are valid check types
		for i := 0; i < len(checksInConfig); i++ {
			if !contains(ohdear.CheckTypes, checksInConfig[i]) {
				return nil, fmt.Errorf("Invalid check type %s - valid check types are 'uptime, 'broken_links', 'mixed_content', 'certificate_health' and 'certificate_transparency'", checksInConfig[i])
			}
		}
		// For each check type specified, see if it is enabled
		for i := 0; i < len(checksWanted); i++ {
			val, ok := checksConfig[checksWanted[i]]

			// If the config specifies that check but not as enabled (true)...
			if ok && !val.(bool) {
				// Delete that check from checks wanted
				checksWanted = append(checksWanted[:i], checksWanted[i+1:]...)
			}
		}
	}

	return checksWanted, nil
}

func resourceOhdearSiteCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("[DEBUG] Calling Create lifecycle function for site")
	checksWanted, err := determineChecksWanted(d)
	if err != nil {
		return fmt.Errorf("Error Creating Site: %s", err.Error())
	}

	site := &ohdear.SiteRequest{
		URL:    d.Get("url").(string),
		TeamID: d.Get("team_id").(int),
		Checks: checksWanted,
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

	err = d.Set("checks", checkStateMapFromSite(newSite))
	if err != nil {
		return fmt.Errorf("Error setting check state: %s", err.Error())
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
				_, err := client.CheckService.DisableCheck(check.ID)
				if err != nil {
					return err
				}
			}
		} else {
			if contains(checks, check.Type) {
				_, err := client.CheckService.EnableCheck(check.ID)
				if err != nil {
					return err
				}
			}
		}
	}

	return resourceOhdearSiteRead(d, meta)
}

func checkStateMapFromSite(site *ohdear.Site) map[string]bool {
	result := make(map[string]bool, len(ohdear.CheckTypes))

	for _, check := range site.Checks {
		result[check.Type] = check.Enabled
	}

	return result
}
