# Mattermost KV Manager Plugin

[![CI](https://github.com/manybugsdev/mattermost-plugin-kv/actions/workflows/ci.yml/badge.svg)](https://github.com/manybugsdev/mattermost-plugin-kv/actions/workflows/ci.yml)

A Mattermost plugin that provides slash commands to perform CRUD (Create, Read, Update, Delete) operations on Mattermost's internal KV (Key-Value) store.

## Features

- **Set**: Create or update a key-value pair
- **Get**: Retrieve the value of a specific key
- **Delete**: Remove a key from the KV store
- **List**: List all keys, with optional prefix filtering

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

The plugin bundle will be created at `dist/com.manybugs.mattermost-plugin-kv-<version>.tar.gz`

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

### Cleaning

```bash
make clean
```

## Usage

After installing and enabling the plugin, you can use the following slash commands:

### Set a key-value pair
```
/kv set <key> <value>
```
Example: `/kv set mykey hello world`

### Get a value by key
```
/kv get <key>
```
Example: `/kv get mykey`

### Delete a key
```
/kv delete <key>
```
Example: `/kv delete mykey`

### List all keys
```
/kv list
```

### List keys with a specific prefix
```
/kv list <prefix>
```
Example: `/kv list my`

### Display help
```
/kv help
```

## Requirements

- Mattermost Server v6.2.1 or later
- Go 1.24 or later (for building from source)

## License

MIT License - see [LICENSE](LICENSE) file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
