# Terraform Provider OhDear

A Terraform Provider for [Oh Dear](https://ohdear.app/).

## Usage

The provider requires an `api_token` (or `OHDEAR_TOKEN` environment variable) and
an optional `team_id` (`OHDEAR_TEAM_ID` environment variable).

<!-- x-release-please-start-version -->
```hcl
terraform {
  required_providers {
    ohdear = {
      source = "articulate/ohdear"
      version = "2.2.1"
    }
  }
}

provider "ohdear" {
  api_token = "my-api-token"
  team_id   = 1234 # optional
}
```
<!-- x-release-please-end -->

To add a site to Oh Dear, create a `ohdear_site` resource.

```hcl
resource "ohdear_site" "test" {
  url = "https://site.iwanttomonitor.com"
}
```

By default, all checks are enabled. You can customize this using the `checks`
block. Any checks not defined in the block are disabled.

```hcl
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

## Contributing

Commit messages must be signed and follow the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/)
format.

## Publishing

Releases are automatically created by [release-please](https://github.com/googleapis/release-please)
on PR merge. This will scan commit messages for new releases based on commit message
and create a release PR. To finish the release, merge the PR, which will kick off
GoReleaser.
