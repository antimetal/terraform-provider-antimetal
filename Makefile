# https://www.gnu.org/prep/standards/standards.html#Makefile-Conventions
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

REGISTRY ?= registry.terraform.io
PROVIDER ?= antimetal
PROJECT_NAME ?= terraform-provider-$(PROVIDER)
VERSION ?= 0.0.1
BINARY_NAME ?= $(PROJECT_NAME)_v$(VERSION)

TOOLSBIN ?= $(shell pwd)/tools/bin
DISTBIN ?= $(shell pwd)/bin

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
ifeq ($(GOARCH),amd64)
GOARCHFLAGS ?= _$(shell go env GOAMD64)
else ifeq ($(GOARCH),arm)
GOARCHFLAGS ?= _$(shell go env GOARM)
else
GOARCHFLAGS ?=
endif

LOCALPLUGINDIR ?= ~/.terraform.d/plugins/$(REGISTRY)/antimetal/$(PROVIDER)/$(VERSION)/$(GOOS)_$(GOARCH)

TESTARGS ?=

## Tools
GOLANGCI_LINT := go run github.com/golangci/golangci-lint/cmd/golangci-lint

GORELEASER         := $(TOOLSBIN)/goreleaser
GORELEASER_VERSION := v1.25.1

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php
.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-22s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

clean: ## Clean bin/ directories.
	rm -rf $(DISTBIN)
	rm -rf $(TOOLSBIN)

##@ Development

.PHONY: test
test: ## Run unit tests.
	go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: testacc
testacc: ## Run Terraform provider acceptance tests.
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: fmt
fmt: ## Run go fmt against code.
	@go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	@go vet ./...

.PHONY: lint
lint: ## Run golangci-lint linters.
	$(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: ## Run golangci-lint linters and perform fixes.
	$(GOLANGCI_LINT) run --fix

.PHONY: generate
generate: ## Run go generate.
	go generate -tags gen ./...

##@ Build

.PHONY: build
build: goreleaser fmt vet ## Build terraform provider using GoReleaser.
	$(GORELEASER) build --snapshot --clean --single-target

.PHONY: install
install: build ## Install provider so that it can be used by Terraform CLI.
	@set -e; { \
		echo "installing provider in $(LOCALPLUGINDIR)" ;\
		mkdir -p $(LOCALPLUGINDIR) ;\
		cp -f $(DISTBIN)/$(PROJECT_NAME)_$(GOOS)_$(GOARCH)$(GOARCHFLAGS)/$(BINARY_NAME) \
			$(LOCALPLUGINDIR)/$(BINARY_NAME) ;\
	}

.PHONY: uninstall
uninstall: ## Uninstall provider.
	rm -f $(LOCALPLUGINDIR)/$(BINARY_NAME)

##@ Tools

# go-install-tool will 'go install' any package with custom target and name of binary, if it doesn't exist
# $1 - target path with name of installed binary
# $2 - package url which can be installed
# $3 - specific version of package
define go-install-tool
@set -e; { \
	binary=$(1)@$(3) ;\
	if [ ! -f $${binary} ]; then \
		package=$(2)@$(3) ;\
		echo "Downloading $${package}" ;\
		GOBIN=$$(dirname $(1)) go install $${package} ;\
		mv $(1) $(1)@$(3) ;\
	fi ;\
}
endef

.PHONY: tools
tools: goreleaser ## Download all tools if neccessary.

.PHONY: goreleaser
goreleaser: $(GORELEASER) ## Download goreleaser locally if neccessary.
$(GORELEASER): $(GORELEASER)@$(GORELEASER_VERSION)
	@ln -sf $< $@
$(GORELEASER)@$(GORELEASER_VERSION):
	$(call go-install-tool,$(GORELEASER),github.com/goreleaser/goreleaser,$(GORELEASER_VERSION))
