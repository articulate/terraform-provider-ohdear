terraform {
  required_providers {
    ohdear = {
      source = "articulate/ohdear"
    }
  }
}

provider "ohdear" {
  api_token = "my-api-token"
}
