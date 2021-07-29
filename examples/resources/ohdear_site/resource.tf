resource "ohdear_site" "test" {
  url = "https://example.com"
}

# Turn off some checks
resource "ohdear_site" "uptime-only" {
  url = "https://example.org"

  broken_links             = false
  certificate_health       = false
  certificate_transparency = false
  mixed_content            = false
}
