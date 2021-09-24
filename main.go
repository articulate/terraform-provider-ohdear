package main

import (
	"context"
	"flag"
	"log"

	"github.com/articulate/terraform-provider-ohdear/internal/provider"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// Format example Terraform files
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

// automatically set by goreleaser
var version = "dev"

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: provider.New(version)}

	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/articulate/ohdear", opts)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}
