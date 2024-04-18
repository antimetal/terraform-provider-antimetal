# https://www.gnu.org/prep/standards/standards.html#Makefile-Conventions
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: build

REGISTRY ?= registry.terraform.io
PROVIDER ?= antimetal
VERSION ?= 0.0.1

LOCALBIN ?= $(shell pwd)/bin
BINARY_NAME ?= terraform-provider-$(PROVIDER)
GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOENV ?= CGO_ENABLED=0 GOBIN=$(LOCALBIN) GOOS=$(GOOS) GOARCH=$(GOARCH)
BUILD_ARGS ?= -ldflags "-X main.Version=$(VERSION)"
TESTARGS ?=

## Tools
GOLANGCI_LINT := github.com/golangci/golangci-lint/cmd/golangci-lint

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

clean: ## Clean bin/ directory.
	rm -rf $(LOCALBIN)

##@ Development

.PHONY: test
test: ## Run tests.
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

.PHONY: fmt
fmt: ## Run go fmt against code.
	@go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	@go vet ./...

.PHONY: lint
lint: ## Run golangci-lint linters.
	go run $(GOLANGCI_LINT) run

.PHONY: lint-fix
lint-fix: ## Run golangci-lint linters and perform fixes.
	go run $(GOLANGCI_LINT) run --fix

.PHONY: generate
generate: ## Run go generate.
	go generate -tags gen ./...

##@ Build

$(LOCALBIN)/$(BINARY_NAME): main.go go.mod go.sum $(shell find internal -type f -name "*.go")
	$(GOENV) go build $(BUILD_ARGS) -o $(LOCALBIN)/$(BINARY_NAME)

.PHONY: build
build: fmt vet $(LOCALBIN)/$(BINARY_NAME) ## Build terraform provider.

.PHONY: install
install: $(LOCALBIN)/$(BINARY_NAME) ## Install provider so that it can be used by Terraform CLI.
	mkdir -p ~/.terraform.d/plugins/$(REGISTRY)/antimetal/$(PROVIDER)/$(VERSION)/$(GOOS)_$(GOARCH)
	ln -sf $(LOCALBIN)/$(BINARY_NAME) ~/.terraform.d/plugins/$(REGISTRY)/antimetal/$(PROVIDER)/$(VERSION)/$(GOOS)_$(GOARCH)/$(BINARY_NAME)_v$(VERSION)
