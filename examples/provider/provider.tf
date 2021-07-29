provider "ohdear" {
  api_token = var.token   # optionally use OHDEAR_TOKEN env var
  api_url   = var.api_url # optionally use OHDEAR_API_URL env var
  team_id   = var.team_id # optionally use OHDEAR_TEAM_ID env var
}
