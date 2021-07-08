# Terraform Provider OhDear

A Terraform Provider for [Oh Dear](https://ohdear.app/).

## Usage

The provider requires an `api_token` (or `OHDEAR_TOKEN` environment variable) and
an optional `api_url` (`OHDEAR_BASE_URL` environment variable, defaults to "https://ohdear.app").

```terraform
provider "ohdear" {
  api_token = "XXXX"
  api_url   = "https://ohdear.app" # optional
}
```

To add a site to Oh Dear, create a `ohdear_site` resource.

```terraform
resource "ohdear_site" "fnord" {
  team_id = 1234
  url     = "https://site.iwanttomonitor.com"
}
```

By default, all site checks are enabled. You can turn off checks by setting them
to false.

```terraform
resource "ohdear_site" "fnord" {
  team_id                  = 1234
  url                      = "https://site.iwanttomonitor.com"
  broken_links             = false
  certificate_health       = false
  certificate_transparency = false
  mixed_content            = false
}
```

## Requirements

* [Terraform](https://www.terraform.io/downloads.html) >= 0.12.x
* [Go](https://golang.org/doc/install) 1.16 (for development)
