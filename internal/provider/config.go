package provider

import (
	"github.com/articulate/terraform-provider-ohdear/pkg/ohdear"
)

type Config struct {
	client *ohdear.Client
	teamID int
}
