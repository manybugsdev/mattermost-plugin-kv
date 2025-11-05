.PHONY: all build clean dist

PLUGIN_ID = com.manybugs.mattermost-plugin-kv
PLUGIN_VERSION = 0.1.0

# Binaries
DIST_DIR = server/dist
BINARY_LINUX = $(DIST_DIR)/plugin-linux-amd64
BINARY_DARWIN = $(DIST_DIR)/plugin-darwin-amd64
BINARY_WINDOWS = $(DIST_DIR)/plugin-windows-amd64.exe

all: clean build

build: $(BINARY_LINUX) $(BINARY_DARWIN) $(BINARY_WINDOWS)

$(BINARY_LINUX):
	mkdir -p $(DIST_DIR)
	cd server && GOOS=linux GOARCH=amd64 go build -o ../$(BINARY_LINUX) .

$(BINARY_DARWIN):
	mkdir -p $(DIST_DIR)
	cd server && GOOS=darwin GOARCH=amd64 go build -o ../$(BINARY_DARWIN) .

$(BINARY_WINDOWS):
	mkdir -p $(DIST_DIR)
	cd server && GOOS=windows GOARCH=amd64 go build -o ../$(BINARY_WINDOWS) .

dist: build
	rm -rf dist
	mkdir -p dist/$(PLUGIN_ID)/server
	cp plugin.json dist/$(PLUGIN_ID)/
	cp -r server/dist dist/$(PLUGIN_ID)/server/
	cd dist && tar -czf $(PLUGIN_ID)-$(PLUGIN_VERSION).tar.gz $(PLUGIN_ID)
	@echo "Plugin package created: dist/$(PLUGIN_ID)-$(PLUGIN_VERSION).tar.gz"

clean:
	rm -rf $(DIST_DIR)
	rm -rf dist
	go clean

deps:
	cd server && go mod download
	cd server && go mod tidy
