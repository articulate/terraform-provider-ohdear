[![Go Report Card](https://goreportcard.com/badge/github.com/articulate/terraform-provider-ohdear)](https://goreportcard.com/report/github.com/articulate/terraform-provider-ohdear)
Terraform Provider OhDear
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

Maintainers
-----------

This provider plugin is maintained by the Terraform team at [Articulate](https://articulate.com/).

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 0.12.x
-	[Go](https://golang.org/doc/install) 1.12 (to build the provider plugin)

Usage
---------------------

This plugin requires two inputs to run: the OhDear organization name and the OhDear api token. The OhDear base url is not required and will default to "OhDear.com" if left out.

You can specify the inputs in your tf plan:

```
provider "ohdear" {
  api_token = "XXXX"
  api_url   = "https://ohdear.app"
}
```

OR you can specify environment variables:

```
OHDEAR_TOKEN=<OhDear api token>
OHDEAR_BASE_URL="https://ohdear.app"
```

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/terraform-providers/terraform-provider-ohdear`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-ohdear
```


```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-ohdear
$ go get -v
$ make build
```

For local development, I've found the below commands helpful. Run them from inside the terraform-provider-ohdear directory

```sh
$ go build -o .terraform/plugins/$GOOS_$GOARCH/terraform-provider-ohdear
$ terraform init -plugin-dir=.terraform/plugins/$GOOS_$GOARCH
```

Using the provider
----------------------

Example terraform plan:

```
provider "ohdear" {
  api_token = "XXXX"
  api_url   = "https://ohdear.app"
}

resource "ohdear_site" "fnord" {
  team_id = 1337
  url     = "https://site.iwanttomonitor.com"
}
```

Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.8+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make bin
...
$ $GOPATH/bin/terraform-provider-ohdear
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

### Best Practices

We are striving to build a provider that is easily consumable and eventually can pass the HashiCorp community audit. In order to achieve this end we must ensure we are following HashiCorp's best practices. This can be derived either from their [documentation on the matter](https://www.terraform.io/docs/extend/best-practices/detecting-drift.html), or by using a simple well written [example as our template](https://github.com/terraform-providers/terraform-provider-datadog).
