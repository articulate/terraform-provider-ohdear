module github.com/articulate/terraform-provider-ohdear

go 1.16

require (
	github.com/articulate/ohdear-sdk v0.0.0-20190816194444-88d39697502e
	github.com/hashicorp/go-cleanhttp v0.5.2
	github.com/hashicorp/terraform-plugin-docs v0.4.0
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.7.0
	gopkg.in/h2non/gock.v1 v1.1.1 // indirect
)

replace github.com/hashicorp/terraform-plugin-docs => github.com/mloberg/terraform-plugin-docs v0.4.1-0.20210603152958-92b22af02b99
