package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/articulate/terraform-provider-ohdear/internal/provider"
	"github.com/articulate/terraform-provider-ohdear/internal/runtime"
)

// Format example Terraform files
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: provider.New,
		ProviderAddr: "registry.terraform.io/articulate/ohdear",
		Debug:        runtime.Debug(),
	})
}
