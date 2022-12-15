package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/articulate/terraform-provider-ohdear/pkg/ohdear"
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
						ohdear.UptimeCheck: {
							Type:        schema.TypeBool,
							Description: "Enable uptime checks.",
							Optional:    true,
						},
						ohdear.BrokenLinksCheck: {
							Type:        schema.TypeBool,
							Description: "Enable broken link checks.",
							Optional:    true,
						},
						ohdear.CertificateHealthCheck: {
							Type:        schema.TypeBool,
							Description: "Enable certificate health checks. Requires the url to use https.",
							Optional:    true,
						},
						ohdear.CertificateTransparencyCheck: {
							Type:        schema.TypeBool,
							Description: "Enable certificate transparency checks. Requires the url to use https.",
							Optional:    true,
						},
						ohdear.MixedContentCheck: {
							Type:        schema.TypeBool,
							Description: "Enable mixed content checks.",
							Optional:    true,
						},
						ohdear.PerformanceCheck: {
							Type:        schema.TypeBool,
							Description: "Enable performance checks.",
							Optional:    true,
						},
						ohdear.DNSCheck: {
							Type:        schema.TypeBool,
							Description: "Enable DNS checks.",
							Default:     false,
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
		return id, fmt.Errorf("corrupted resource ID in terraform state, Oh Dear only supports integer IDs. Err: %w", err)
	}
	return id, nil
}

func resourceOhdearSiteDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	checks := d.Get("checks").([]interface{})
	if len(checks) == 0 {
		isHTTPS := strings.HasPrefix(d.Get("url").(string), "https")
		checks = append(checks, map[string]bool{
			ohdear.UptimeCheck:                  true,
			ohdear.BrokenLinksCheck:             true,
			ohdear.CertificateHealthCheck:       isHTTPS,
			ohdear.CertificateTransparencyCheck: isHTTPS,
			ohdear.MixedContentCheck:            isHTTPS,
			ohdear.PerformanceCheck:             true,
			ohdear.DNSCheck:                     false, // TODO: turn to true on next major release (breaking change)
		})

		if err := d.SetNew("checks", checks); err != nil {
			return fmt.Errorf("could not set checks: %w", err)
		}
	}

	// set team_id from provider default if not provided
	if d.Get("team_id") != 0 {
		return nil
	}

	if err := d.SetNew("team_id", meta.(*Config).teamID); err != nil {
		return fmt.Errorf("could not set team_id: %w", err)
	}

	return nil
}

func resourceOhdearSiteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Calling Create lifecycle function for site")

	client := meta.(*Config).client
	site, err := client.AddSite(d.Get("url").(string), d.Get("team_id").(int), checksWanted(d))
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
	site, err := client.GetSite(id)
	if err != nil {
		return diagErrorf(err, "Could not find site %d in Oh Dear", id)
	}

	checks := checkStateMapFromSite(site)
	if err := d.Set("checks", []interface{}{checks}); err != nil {
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

	client := meta.(*Config).client
	if err = client.RemoveSite(id); err != nil {
		return diagErrorf(err, "Could not remove site %d from Oh Dear", id)
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
	site, err := client.GetSite(id)
	if err != nil {
		return diagErrorf(err, "Could not find site in Oh Dear")
	}

	// Sync downstream checks with config
	checksWanted := checksWanted(d)
	for _, check := range site.Checks {
		if check.Enabled {
			if !contains(checksWanted, check.Type) {
				if err := client.DisableCheck(check.ID); err != nil {
					return diagErrorf(err, "Could not remove check to site in Oh Dear")
				}
			}
		} else {
			if contains(checksWanted, check.Type) {
				if err := client.EnableCheck(check.ID); err != nil {
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
		if contains(ohdear.AllChecks, check.Type) {
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
