package main

import (
	"github.com/articulate/terraform-provider-ohdear/ohdear"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{ProviderFunc: ohdear.Provider})
}
