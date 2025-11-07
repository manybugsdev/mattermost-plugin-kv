# Mattermost Hello World Plugin

[![CI](https://github.com/manybugsdev/mattermost-plugin-kv/actions/workflows/ci.yml/badge.svg)](https://github.com/manybugsdev/mattermost-plugin-kv/actions/workflows/ci.yml)

A simple Mattermost plugin that demonstrates basic plugin functionality. This plugin serves as a "Hello World" example following the [Mattermost plugin development documentation](https://developers.mattermost.com/integrate/plugins/components/server/hello-world/).

## Features

This plugin demonstrates two key plugin hooks:

- **OnActivate**: Logs a message when the plugin is activated
- **MessageHasBeenPosted**: Responds to every message posted in Mattermost by:
  - Logging the message content to the server logs
  - Sending an ephemeral "Hello from Hello World plugin!" message back to the user

## Installation

1. Download the latest release from the releases page
2. Upload the plugin to your Mattermost server via System Console → Plugins → Plugin Management
3. Enable the plugin

## Building from Source

This plugin uses the [Mattermost Plugin Starter Template](https://github.com/mattermost/mattermost-plugin-starter-template) structure.

### Prerequisites

- Go 1.24 or later
- Make

### Build

```bash
make
```

To create a distributable package:

```bash
make dist
```

The plugin bundle will be created at `dist/com.manybugs.mattermost-plugin-hello-world-<version>.tar.gz`

## Development

### Building

```bash
# Install dependencies
make deps

# Build for all platforms
make server

# Build and bundle
make dist
```

### Testing

```bash
make test
```

### Cleaning

```bash
make clean
```

## Usage

Once the plugin is activated, it will:
1. Log "Hello World plugin has been activated!" to the Mattermost server logs
2. Respond to every message posted by any user with an ephemeral message visible only to that user

To see the plugin in action:
- Post any message in any channel
- You will receive an ephemeral message: "Hello from Hello World plugin!"

## Requirements

- Mattermost Server v6.2.1 or later
- Go 1.24 or later (for building from source)

## License

MIT License - see [LICENSE](LICENSE) file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Reference

This plugin is based on the [Mattermost Plugin Hello World Guide](https://developers.mattermost.com/integrate/plugins/components/server/hello-world/).
