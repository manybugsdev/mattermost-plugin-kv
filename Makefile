# Mattermost Plugin Makefile
# Simplified version for server-only plugins

GO ?= $(shell command -v go 2> /dev/null)
GOPATH ?= $(shell go env GOPATH)
GOBIN ?= $(shell go env GOPATH)/bin

# Plugin configuration
PLUGIN_ID = com.manybugs.mattermost-plugin-kv
PLUGIN_VERSION = $(shell git describe --tags --always --dirty 2>/dev/null || echo "0.0.0")
BUNDLE_NAME = $(PLUGIN_ID)-$(PLUGIN_VERSION).tar.gz

# Directories
DIST_DIR = dist
SERVER_DIST_DIR = server/dist
ASSETS_DIR = assets

# Build flags
GO_BUILD_FLAGS ?=
GO_TEST_FLAGS ?= -race

export GO111MODULE=on

.PHONY: default
default: all

.PHONY: all
all: check-style test dist

## Build server binaries for all platforms
.PHONY: server
server:
	@echo Building plugin server
	@mkdir -p $(SERVER_DIST_DIR)
	cd server && env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -trimpath -o dist/plugin-linux-amd64
	cd server && env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GO) build $(GO_BUILD_FLAGS) -trimpath -o dist/plugin-linux-arm64
	cd server && env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -trimpath -o dist/plugin-darwin-amd64
	cd server && env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 $(GO) build $(GO_BUILD_FLAGS) -trimpath -o dist/plugin-darwin-arm64
	cd server && env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GO) build $(GO_BUILD_FLAGS) -trimpath -o dist/plugin-windows-amd64.exe

## Create plugin bundle
.PHONY: bundle
bundle:
	@echo Creating plugin bundle
	@rm -rf $(DIST_DIR)
	@mkdir -p $(DIST_DIR)/$(PLUGIN_ID)/server
	@cp -r $(SERVER_DIST_DIR) $(DIST_DIR)/$(PLUGIN_ID)/server/
	@cp plugin.json $(DIST_DIR)/$(PLUGIN_ID)/
	@if [ -d "$(ASSETS_DIR)" ]; then cp -r $(ASSETS_DIR) $(DIST_DIR)/$(PLUGIN_ID)/; fi
	@cd $(DIST_DIR) && tar -czf $(BUNDLE_NAME) $(PLUGIN_ID)
	@echo Plugin bundle created at: $(DIST_DIR)/$(BUNDLE_NAME)

## Build and bundle the plugin
.PHONY: dist
dist: server bundle

## Run tests
.PHONY: test
test:
	@echo Running tests
	@$(GO) test $(GO_TEST_FLAGS) -v ./...

## Install Go tools for linting
.PHONY: install-go-tools
install-go-tools:
	@echo Installing Go tools
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.8

## Check code style
.PHONY: check-style
check-style: install-go-tools
	@echo Checking code style
	@$(GO) vet ./...
	@$(GOBIN)/golangci-lint run ./...

## Download and tidy Go dependencies
.PHONY: deps
deps:
	@$(GO) mod download
	@$(GO) mod tidy

## Clean build artifacts
.PHONY: clean
clean:
	@echo Cleaning build artifacts
	@rm -rf $(DIST_DIR)
	@rm -rf $(SERVER_DIST_DIR)
	@rm -rf server/coverage.txt

## Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all           - Run check-style, test, and dist (default)"
	@echo "  server        - Build server binaries for all platforms"
	@echo "  bundle        - Create plugin bundle"
	@echo "  dist          - Build and bundle the plugin"
	@echo "  test          - Run tests"
	@echo "  check-style   - Run linters and code checks"
	@echo "  deps          - Download and tidy Go dependencies"
	@echo "  clean         - Remove build artifacts"
	@echo "  help          - Show this help message"
