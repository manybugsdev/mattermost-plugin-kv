# Mattermost KV Manager Plugin

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

```bash
make build
```

To create a distributable package:

```bash
make dist
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

- Mattermost Server v5.20.0 or later
- Go 1.21 or later (for building from source)

## License

MIT License - see [LICENSE](LICENSE) file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
