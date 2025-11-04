# Default version for local builds
VERSION ?= 0.1.0

# Default install path based on OS and architecture
OS_ARCH := $(shell go env GOOS)_$(shell go env GOARCH)
INSTALL_PATH := ~/.terraform.d/plugins/registry.terraform.io/hqarroum/paddle/$(VERSION)/$(OS_ARCH)

## Show this help
.PHONY: help
help:
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'

## Build the provider
.PHONY: build
build:
	go build -o terraform-provider-paddle

## Install the provider locally for development
.PHONY: install
install: build
	@echo "Installing provider to $(INSTALL_PATH)"
	@mkdir -p $(INSTALL_PATH)
	@cp terraform-provider-paddle $(INSTALL_PATH)/
	@echo "Provider installed successfully!"

## Run unit tests
.PHONY: test
test:
	go test -v -cover ./...

## Run acceptance tests (requires PADDLE_API_KEY)
.PHONY: testacc
testacc:
	TF_ACC=1 go test -v -cover -timeout 120m ./...

## Format Go code
.PHONY: fmt
fmt:
	go fmt ./...

## Run go vet
.PHONY: vet
vet:
	go vet ./...

## Run golangci-lint
.PHONY: lint
lint:
	golangci-lint run

## Clean build artifacts
.PHONY: clean
clean:
	rm -f terraform-provider-paddle
	go clean -cache

## Generate documentation
.PHONY: docs
docs:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate

## Download dependencies
.PHONY: deps
deps:
	go mod download
	go mod tidy

## Build and run provider in debug mode
.PHONY: debug
debug: build
	@echo "Starting provider in debug mode..."
	@echo "Set TF_REATTACH_PROVIDERS environment variable from the output"
	./terraform-provider-paddle --debug

.PHONY: all
all: fmt vet build test
