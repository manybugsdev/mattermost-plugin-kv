# Mattermost KV Store Manager Plugin

A plugin for managing key-value pairs in Mattermost's built-in KV database.

This plugin provides a `/kv` slash command that allows you to perform CRUD operations on the plugin's KV store directly from Mattermost.

## Features

- **Set** key-value pairs
- **Get** values by key
- **Delete** individual keys
- **List** all stored keys
- **Delete all** keys at once
- User-friendly command interface with help documentation

## Usage

Once the plugin is installed and activated, you can use the `/kv` command in any channel:

### Available Commands

- `/kv set <key> <value>` - Set a key-value pair
- `/kv get <key>` - Get the value for a key
- `/kv delete <key>` - Delete a key-value pair
- `/kv list` - List all keys in the store
- `/kv deleteall` - Delete all key-value pairs
- `/kv help` - Show help message

### Examples

```
/kv set mykey Hello World
/kv get mykey
/kv list
/kv delete mykey
/kv deleteall
```

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

