# Terraform Provider OhDear

A Terraform Provider for [Oh Dear](https://ohdear.app/).

## Usage

The provider requires an `api_token` (or `OHDEAR_TOKEN` environment variable) and
an optional `team_id` (`OHDEAR_TEAM_ID` environment variable).

```terraform
provider "ohdear" {
  api_token = "my-api-token"
  team_id   = 1234 # optional
}
```

To add a site to Oh Dear, create a `ohdear_site` resource.

```terraform
resource "ohdear_site" "test" {
  url = "https://site.iwanttomonitor.com"
}
```

By default, all checks are enabled. You can customize this using the `checks`
block. Any checks not defined in the block are disabled.

```terraform
resource "ohdear_site" "test" {
  url = "https://site.iwanttomonitor.com"

  checks {
    uptime = true
  }
}
```

## Development Requirements

* [Go](https://golang.org/doc/install) (for development)
* [golangci-lint](https://golangci-lint.run/)
* [GoReleaser](https://goreleaser.com/)
