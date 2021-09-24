resource "ohdear_site" "test" {
  url = "https://example.com"
  # all checks are enabled
}

resource "ohdear_site" "uptime-only" {
  url = "https://example.org"

  # Only the uptime check is enabled
  checks {
    uptime = true
  }
}
