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

func resourceOhdearMonitor() *schema.Resource {
	return &schema.Resource{
		Description:   "`ohdear_monitor` manages a monitor in Oh Dear.",
		CreateContext: resourceOhdearMonitorCreate,
		ReadContext:   resourceOhdearMonitorRead,
		DeleteContext: resourceOhdearMonitorDelete,
		UpdateContext: resourceOhdearMonitorUpdate,
		Schema: map[string]*schema.Schema{
			"url": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "URL of the monitor to be checked.",
				ValidateFunc: validation.IsURLWithHTTPorHTTPS,
			},
			"team_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "ID of the team for this monitor. If not set, will use `team_id` configured in provider.",
			},
			"checks": {
				Type:        schema.TypeList,
				Description: "Set the checks enabled for the monitor. If block is not present, it will enable all checks.",
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
							Deprecated: "This check was removed by OhDear and will be removed " +
								"in a future major release.",
							DiffSuppressFunc: func(_, _, _ string, _ *schema.ResourceData) bool {
								return true
							},
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
							Deprecated: "This check was merged with the 'uptime' check by OhDear and will be removed " +
								"in a future major release.",
							DiffSuppressFunc: func(_, _, _ string, _ *schema.ResourceData) bool {
								return true
							},
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
		CustomizeDiff: resourceOhdearMonitorDiff,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func getMonitorID(d *schema.ResourceData) (int, error) {
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return id, fmt.Errorf("corrupted resource ID in terraform state, Oh Dear only supports integer IDs. Err: %w", err)
	}
	return id, err
}

func resourceOhdearMonitorDiff(_ context.Context, d *schema.ResourceDiff, meta interface{}) error {
	checks := d.Get("checks").([]interface{})
	if len(checks) == 0 {
		isHTTPS := strings.HasPrefix(d.Get("url").(string), "https")
		checks = append(checks, map[string]bool{
			ohdear.UptimeCheck:            true,
			ohdear.BrokenLinksCheck:       true,
			ohdear.CertificateHealthCheck: isHTTPS,
			ohdear.MixedContentCheck:      isHTTPS,
			ohdear.DNSCheck:               false, // TODO: turn to true on next major release (breaking change)
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

func resourceOhdearMonitorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Println("[DEBUG] Calling Create lifecycle function for monitor")

	client := meta.(*Config).client
	monitor, err := client.AddMonitor(d.Get("url").(string), d.Get("team_id").(int), checksWanted(d))
	if err != nil {
		return diagErrorf(err, "Could not add monitor to Oh Dear")
	}

	d.SetId(strconv.Itoa(monitor.ID))

	return resourceOhdearMonitorRead(ctx, d, meta)
}

func resourceOhdearMonitorRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Calling Read lifecycle function for monitor %s\n", d.Id())

	id, err := getMonitorID(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*Config).client
	monitor, err := client.GetMonitor(id)
	if err != nil {
		return diagErrorf(err, "Could not find monitor %d in Oh Dear", id)
	}

	checks := checkStateMapFromMonitor(monitor)
	if err := d.Set("checks", []interface{}{checks}); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("url", monitor.URL); err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("team_id", monitor.TeamID); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceOhdearMonitorDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Calling Delete lifecycle function for monitor %s\n", d.Id())

	id, err := getMonitorID(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*Config).client
	if err = client.RemoveMonitor(id); err != nil {
		return diagErrorf(err, "Could not remove monitor %d from Oh Dear", id)
	}

	return nil
}

func resourceOhdearMonitorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	log.Printf("[DEBUG] Calling Update lifecycle function for monitor %s\n", d.Id())

	id, err := getMonitorID(d)
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*Config).client
	monitor, err := client.GetMonitor(id)
	if err != nil {
		return diagErrorf(err, "Could not find monitor in Oh Dear")
	}

	// Sync downstream checks with config
	checksWanted := checksWanted(d)
	for _, check := range monitor.Checks {
		if check.Enabled {
			if !contains(checksWanted, check.Type) {
				if err := client.DisableCheck(check.ID); err != nil {
					return diagErrorf(err, "Could not remove check to monitor in Oh Dear")
				}
			}
		} else {
			if contains(checksWanted, check.Type) {
				if err := client.EnableCheck(check.ID); err != nil {
					return diagErrorf(err, "Could not add check to monitor in Oh Dear")
				}
			}
		}
	}

	return resourceOhdearMonitorRead(ctx, d, meta)
}

func checkStateMapFromMonitor(monitor *ohdear.Monitor) map[string]bool {
	result := make(map[string]bool)
	for _, check := range monitor.Checks {
		if contains(ohdear.AllChecks, check.Type) {
			result[check.Type] = check.Enabled
		}
	}

	return result
}
