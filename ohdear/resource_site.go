package ohdear

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/articulate/ohdear-sdk/ohdear"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceOhdearSite() *schema.Resource {
	return &schema.Resource{
		Description: "`ohdear_site` manages a site in Oh Dear.",

		CreateContext: resourceOhdearSiteCreate,
		ReadContext:   resourceOhdearSiteRead,
		DeleteContext: resourceOhdearSiteDelete,
		UpdateContext: resourceOhdearSiteUpdate,
		Schema: map[string]*schema.Schema{
			"url": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "URL of the site to be checked.",
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"team_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "ID of the team for this site. If not set, will use `team_id` configured in provider.",
			},
			"checks": {
				Type:        schema.TypeList,
				Description: "Set the checks enabled for the site. If block is not present, it will enable all checks.",
				Optional:    true,
				Computed:    true,
				MaxItems:    1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"uptime": {
							Type:        schema.TypeBool,
							Description: "Enable uptime checks.",
							Optional:    true,
						},
						"broken_links": {
							Type:        schema.TypeBool,
							Description: "Enable broken link checks.",
							Optional:    true,
						},
						"certificate_health": {
							Type:        schema.TypeBool,
							Description: "Enable certificate health checks. Requires the url to use https.",
							Optional:    true,
						},
						"certificate_transparency": {
							Type:        schema.TypeBool,
							Description: "Enable certificate transparency checks. Requires the url to use https.",
							Optional:    true,
						},
						"mixed_content": {
							Type:        schema.TypeBool,
							Description: "Enable mixed content checks.",
							Optional:    true,
						},
						"performance": {
							Type:        schema.TypeBool,
							Description: "Enable performance checks.",
							Optional:    true,
						},
					},
				},
			},
		},
		CustomizeDiff: resourceOhdearSiteDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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

func resourceOhdearSiteDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	checks := d.Get("checks").([]interface{})
	if len(checks) == 0 {
		isHTTPS := strings.HasPrefix(d.Get("url").(string), "https")
		checks = append(checks, map[string]bool{
			"uptime":                   true,
			"broken_links":             true,
			"certificate_health":       isHTTPS,
			"certificate_transparency": isHTTPS,
			"mixed_content":            true,
			"performance":              true,
		})

		if err := d.SetNew("checks", checks); err != nil {
			return err
		}
	}

	// set team_id from provider default if not provided
	if d.Get("team_id") == 0 {
		return d.SetNew("team_id", meta.(*Config).teamID)
	}

	return nil
}

func resourceOhdearSiteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Calling Create lifecycle function for site")

	client := meta.(*Config).client
	site, _, err := client.SiteService.CreateSite(&ohdear.SiteRequest{
		URL:    d.Get("url").(string),
		TeamID: d.Get("team_id").(int),
		Checks: checksWanted(d),
	})

	if err != nil {
		return diagErrorf(err, "Could not add site to Oh Dear")
	}

	d.SetId(strconv.Itoa(site.ID))

	return resourceOhdearSiteRead(ctx, d, meta)
}

func resourceOhdearSiteRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Calling Read lifecycle function for site %s\n", d.Id())

	id, err := getSiteID(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*Config).client
	site, _, err := client.SiteService.GetSite(id)
	if err != nil {
		return diagErrorf(err, "Could not find site in Oh Dear")
	}

	checkStateMap := checkStateMapFromSite(site)
	if err := d.Set("checks", []interface{}{checkStateMap}); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("url", site.URL); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("team_id", site.TeamID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceOhdearSiteDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Calling Delete lifecycle function for site %s\n", d.Id())

	id, err := getSiteID(d)
	if err != nil {
		return diag.FromErr(err)
	}

	if _, err = meta.(*Config).client.SiteService.DeleteSite(id); err != nil {
		return diagErrorf(err, "Could not remove site from Oh Dear")
	}

	return nil
}

func resourceOhdearSiteUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Calling Update lifecycle function for site %s\n", d.Id())

	id, err := getSiteID(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*Config).client
	site, _, err := client.SiteService.GetSite(id)
	if err != nil {
		return diagErrorf(err, "Could not find site in Oh Dear")
	}

	// Sync downstream checks with config
	checksWanted := checksWanted(d)
	for _, check := range site.Checks {
		if check.Enabled {
			if !contains(checksWanted, check.Type) {
				if _, err := client.CheckService.DisableCheck(check.ID); err != nil {
					return diagErrorf(err, "Could not remove check to site in Oh Dear")
				}
			}
		} else {
			if contains(checksWanted, check.Type) {
				if _, err := client.CheckService.EnableCheck(check.ID); err != nil {
					return diagErrorf(err, "Could not add check to site in Oh Dear")
				}
			}
		}
	}

	return resourceOhdearSiteRead(ctx, d, meta)
}

func checkStateMapFromSite(site *ohdear.Site) map[string]bool {
	result := make(map[string]bool)
	for _, check := range site.Checks {
		if check.Type != "cron" {
			result[check.Type] = check.Enabled
		}
	}

	return result
}

func checksWanted(d *schema.ResourceData) []string {
	checks := []string{}
	schema := d.Get("checks").([]interface{})[0].(map[string]interface{})
	for check, enabled := range schema {
		if enabled.(bool) {
			checks = append(checks, check)
		}
	}

	return checks
}
