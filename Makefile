# Makefile for Mattermost plugin

# Plugin configuration
PLUGIN_ID := $(shell cat plugin.json | grep '"id"' | sed 's/.*"id": "\(.*\)",/\1/')
PLUGIN_VERSION := $(shell cat plugin.json | grep '"version"' | sed 's/.*"version": "\(.*\)",/\1/')

# Build configuration
PLUGIN_EXECUTABLE := plugin
PLUGIN_ARCHIVE := plugin.tar.gz

# Go build flags
GO_BUILD_FLAGS := -trimpath

.PHONY: all
all: $(PLUGIN_EXECUTABLE) $(PLUGIN_ARCHIVE)

.PHONY: clean
clean:
	rm -f $(PLUGIN_EXECUTABLE) $(PLUGIN_ARCHIVE)

$(PLUGIN_EXECUTABLE): plugin.go go.mod go.sum
	go build $(GO_BUILD_FLAGS) -o $(PLUGIN_EXECUTABLE) plugin.go

$(PLUGIN_ARCHIVE): $(PLUGIN_EXECUTABLE) plugin.json
	tar -czvf $(PLUGIN_ARCHIVE) plugin.json $(PLUGIN_EXECUTABLE)

.PHONY: check-style
check-style:
	@if ! command -v golangci-lint > /dev/null; then \
		echo "golangci-lint is not installed. Please install it to run linting."; \
		exit 1; \
	fi
	golangci-lint run ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: help
help:
	@echo "Mattermost Plugin Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  all           - Build plugin.exe and plugin.tar.gz (default)"
	@echo "  clean         - Remove built files"
	@echo "  check-style   - Run linting checks"
	@echo "  test          - Run tests"
	@echo "  help          - Show this help message"
