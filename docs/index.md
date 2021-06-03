---
page_title: "Oh Dear Provider"
subcategory: ""
description: |-

---

# Oh Dear Provider

Setup monitors with [Oh Dear](https://ohdear.app/).

## Example Usage

```terraform
provider "ohdear" {
  api_token = "my-api-token"
}

resource "ohdear_site" "fnord" {
  team_id = 1337
  url     = "https://site.iwanttomonitor.com"
}
```

## Schema

### Required

- **api_token** (String) or via environment variable `OHDEAR_TOKEN`

### Optional

- **api_url** (String) or via environment variable `OHDEAR_BASE_URL` (default: _https://ohdear.app_)
