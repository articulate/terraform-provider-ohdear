package main

import (
	"context"
	"flag"
	"log"

	"github.com/articulate/terraform-provider-ohdear/ohdear"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

// automatically set by goreleaser
var version = "dev"

func main() {
	var debugMode bool

	flag.BoolVar(&debugMode, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := &plugin.ServeOpts{ProviderFunc: ohdear.New(version)}

	if debugMode {
		err := plugin.Debug(context.Background(), "registry.terraform.io/articulate/ohdear", opts)
		if err != nil {
			log.Fatal(err.Error())
		}
		return
	}

	plugin.Serve(opts)
}
