# Mattermost Hello World Plugin

A simple hello-world plugin for Mattermost.

This plugin demonstrates the basic structure of a Mattermost plugin. When a user posts "Hello, world!", the plugin will modify the message by appending a note that it was modified by the plugin.

## Features

- Intercepts messages containing "Hello, world!"
- Modifies the message to demonstrate plugin functionality

## Building

To build the plugin using the Makefile:

```bash
make
```

This will create both `plugin.exe` and `plugin.tar.gz`.

To clean build artifacts:

```bash
make clean
```

For help with available make targets:

```bash
make help
```

## Installation

1. Build the plugin
2. Create a plugin bundle with the compiled binary and plugin.json
3. Upload to your Mattermost server via System Console → Plugins → Plugin Management

## License

MIT License - see [LICENSE](LICENSE) file for details

