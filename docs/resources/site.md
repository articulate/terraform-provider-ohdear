---
page_title: "ohdear_site Resource - terraform-provider-ohdear"
subcategory: ""
description: |-

---

# ohdear_site (Resource)

Manages a monitored site in Oh Dear.

## Example Usage

```terraform
resource "ohdear_site" "example" {
  team_id      = 1337
  url          = "https://example.com"
  broken_links = false
}
```

## Schema

### Required

- **team_id** (Number) ID of the team for this site
- **url** (String) URL of the site to be checked

### Optional

- **broken_links** (Boolean) Enable/Disable broken_links check (default: *true*)
- **certificate_health** (Boolean) Enable/Disable certificate_health check (default: *true*)
- **certificate_transparency** (Boolean) Enable/Disable certificate_transparency check. Cannot be used with http URLs (default: *true*)
- **id** (String) The ID of this resource.
- **mixed_content** (Boolean) Enable/Disable mixed_content check (default: *true*)
- **uptime** (Boolean) Enable/Disable uptime check (default: *true*)
