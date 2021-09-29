PKG_LIST := $(shell go list ./... | grep -v /vendor/)
OS_ARCH=$(shell go env GOOS)_$(shell go env GOARCH)
HOSTNAME=registry.terraform.io
NAMESPACE=articulate
NAME=ohdear
VERSION=$(shell git describe --abbrev=0 | sed 's/^v//')

help:
	@echo "+ $@"
	@grep -hE '(^[a-zA-Z0-9\._-]+:.*?##.*$$)|(^##)' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}{printf "\033[32m%-30s\033[0m %s\n", $$1, $$2}' | sed -e 's/\[32m## /[33m/'
.PHONY: help

##
## Build
## ---------------------------------------------------------------------------

build: ## Build for current OS/Arch
	@echo "+ $@"
	@goreleaser build --rm-dist --skip-validate --single-target
.PHONY: build

all: ## Build all OS/Arch
	@echo "+ $@"
	@goreleaser build --rm-dist --skip-validate
.PHONY: all

install: build ## Install to global Terraform plugin directory
	@echo "+ $@"
	@mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	@mv dist/terraform-provider-${NAME}_${OS_ARCH}/terraform-provider-* ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
.PHONY: install

generate: ## Autogenerate docs and resources
	@echo "+ $@"
	@go generate ${PKG_LIST}
.PHONY: generate

##
## Development
## ---------------------------------------------------------------------------

dev: ## Start development environment via Docker
	@echo "+ $@"
	@docker build -t ${NAMESPACE}/terraform-provider-${NAME} .
	@docker run --rm -it -v $(PWD):/go/src/github.com/${NAMESPACE}/terraform-provider-${NAME} -w /go/src/github.com/${NAMESPACE}/terraform-provider-${NAME} ${NAMESPACE}/terraform-provider-${NAME}
.PHONY: dev

mod: ## Make sure go.mod is up to date
	@echo "+ $@"
	@go mod tidy
.PHONY: mod

lint: ## Lint Go code
	@echo "+ $@"
	@golangci-lint run
.PHONY: lint

fix: ## Try to fix lint issues
	@echo "+ $@"
	@golangci-lint run --fix
.PHONY: fix

##
## Tests
## ---------------------------------------------------------------------------

test: ## Run tests
	@echo "+ $@"
	@go test ${PKG_LIST} -v $(TESTARGS) -parallel=4
.PHONY: test

testacc: ## Run acceptance tests
	@echo "+ $@"
	@TF_ACC=1 go test ${PKG_LIST} -v -cover $(TESTARGS) -timeout 120m
.PHONY: testacc


# Print the value of any variable as make print-VAR
print-%  : ; @echo $* = $($*)
