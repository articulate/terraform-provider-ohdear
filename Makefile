PKG_LIST := $(shell go list ./... | grep -v /vendor/)

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

##
## Development
## ---------------------------------------------------------------------------

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
