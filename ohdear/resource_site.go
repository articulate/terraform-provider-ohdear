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
			"uptime": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable/Disable uptime check",
				Default:     true,
			},
			"broken_links": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable/Disable broken_links check",
				Default:     true,
			},
			"certificate_health": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable/Disable certificate_health check",
				Default:     true,
			},
			"mixed_content": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable/Disable mixed_content check",
				Default:     true,
			},
			"certificate_transparency": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Enable/Disable certificate_transparency check. Cannot be used with http URLs",
				Default:     true,
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

	checkStateMap := checkStateMapFromSite(newSite)
	for _, checkType := range ohdear.CheckTypes {
		err := d.Set(checkType, checkStateMap[checkType])
		if err != nil {
			return fmt.Errorf("Error setting check state: %s", err.Error())
		}
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

// determineChecksWanted returns the types of checks which are specified
// explicitly as enabled in our config OR which are implicitly enabled
// by virtue of their exclusion from config.
func determineChecksWanted(d *schema.ResourceData) ([]string, error) {
	// We want all the checks by default...
	checksWanted := make([]string, len(ohdear.CheckTypes))

	for _, checkType := range ohdear.CheckTypes {
		config := d.Get(checkType).(bool)
		if config {
			checksWanted = append(checksWanted, checkType)
		}
	}

	return checksWanted, nil
}
