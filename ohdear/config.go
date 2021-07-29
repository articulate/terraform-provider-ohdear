package ohdear

import "github.com/articulate/ohdear-sdk/ohdear"

type Config struct {
	apiToken string
	baseURL  string
	teamID   int

	client *ohdear.Client
}
