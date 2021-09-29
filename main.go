package main

import (
	"context"
	"log"

	"github.com/articulate/terraform-provider-ohdear/internal/provider"
	"github.com/articulate/terraform-provider-ohdear/internal/runtime"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// Format example Terraform files
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	opts := &plugin.ServeOpts{ProviderFunc: provider.New}

	if runtime.Debug() {
		err := plugin.Debug(context.Background(), "registry.terraform.io/articulate/ohdear", opts)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}
