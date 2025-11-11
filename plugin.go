package main

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin
}

// OnActivate is invoked when the plugin is activated.
func (p *Plugin) OnActivate() error {
	// Register the /kv command
	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          "kv",
		AutoComplete:     true,
		AutoCompleteDesc: "Manage key-value pairs in the plugin KV store",
		AutoCompleteHint: "[set|get|delete|list|deleteall|help]",
		DisplayName:      "KV Store Management",
		Description:      "CRUD operations for the plugin KV store",
	}); err != nil {
		return fmt.Errorf("failed to register command: %w", err)
	}

	return nil
}

// ExecuteCommand executes the /kv command
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	// Parse the command
	split := strings.Fields(args.Command)
	if len(split) < 2 {
		return p.getHelpResponse(), nil
	}

	action := split[1]

	switch action {
	case "set":
		return p.handleSet(split)
	case "get":
		return p.handleGet(split)
	case "delete":
		return p.handleDelete(split)
	case "list":
		return p.handleList(split)
	case "deleteall":
		return p.handleDeleteAll()
	case "help":
		return p.getHelpResponse(), nil
	default:
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Unknown action: %s\nUse `/kv help` for usage information.", action),
		}, nil
	}
}

// handleSet stores a key-value pair
func (p *Plugin) handleSet(args []string) (*model.CommandResponse, *model.AppError) {
	if len(args) < 4 {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Usage: `/kv set <key> <value>`",
		}, nil
	}

	key := args[2]
	value := strings.Join(args[3:], " ")

	if err := p.API.KVSet(key, []byte(value)); err != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Error setting key: %s", err.Error()),
		}, nil
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         fmt.Sprintf("✓ Key `%s` set successfully", key),
	}, nil
}

// handleGet retrieves a value by key
func (p *Plugin) handleGet(args []string) (*model.CommandResponse, *model.AppError) {
	if len(args) < 3 {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Usage: `/kv get <key>` or `/kv get <pluginid:key>`",
		}, nil
	}

	keyArg := args[2]

	// Check if the key is in the format pluginid:key
	if strings.Contains(keyArg, ":") {
		parts := strings.SplitN(keyArg, ":", 2)
		if len(parts) == 2 {
			pluginID := parts[0]
			key := parts[1]

			value, err := p.getPluginKVValue(pluginID, key)
			if err != nil {
				return &model.CommandResponse{
					ResponseType: model.CommandResponseTypeEphemeral,
					Text:         fmt.Sprintf("Error getting key from plugin `%s`: %s", pluginID, err.Error()),
				}, nil
			}

			return &model.CommandResponse{
				ResponseType: model.CommandResponseTypeEphemeral,
				Text:         fmt.Sprintf("**Plugin:** `%s`\n**Key:** `%s`\n**Value:** %s", pluginID, key, string(value)),
			}, nil
		}
	}

	// Default behavior: get from this plugin's KV store
	key := keyArg
	value, err := p.API.KVGet(key)
	if err != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Error getting key: %s", err.Error()),
		}, nil
	}

	if value == nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Key `%s` not found", key),
		}, nil
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         fmt.Sprintf("**Key:** `%s`\n**Value:** %s", key, string(value)),
	}, nil
}

// handleDelete removes a key-value pair
func (p *Plugin) handleDelete(args []string) (*model.CommandResponse, *model.AppError) {
	if len(args) < 3 {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "Usage: `/kv delete <key>`",
		}, nil
	}

	key := args[2]

	if err := p.API.KVDelete(key); err != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Error deleting key: %s", err.Error()),
		}, nil
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         fmt.Sprintf("✓ Key `%s` deleted successfully", key),
	}, nil
}

// handleList lists all keys
func (p *Plugin) handleList(args []string) (*model.CommandResponse, *model.AppError) {
	// Check if --all flag is present
	showAll := false
	for _, arg := range args {
		if arg == "--all" {
			showAll = true
			break
		}
	}

	if showAll {
		return p.handleListAll()
	}

	page := 0
	perPage := 100

	keys, err := p.API.KVList(page, perPage)
	if err != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Error listing keys: %s", err.Error()),
		}, nil
	}

	if len(keys) == 0 {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "No keys found in the KV store",
		}, nil
	}

	var text strings.Builder
	text.WriteString(fmt.Sprintf("**KV Store Keys** (showing %d keys):\n", len(keys)))
	for i, key := range keys {
		text.WriteString(fmt.Sprintf("%d. `%s`\n", i+1, key))
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         text.String(),
	}, nil
}

// handleDeleteAll removes all key-value pairs
func (p *Plugin) handleDeleteAll() (*model.CommandResponse, *model.AppError) {
	if err := p.API.KVDeleteAll(); err != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Error deleting all keys: %s", err.Error()),
		}, nil
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         "✓ All keys deleted successfully",
	}, nil
}

// getHelpResponse returns help information
func (p *Plugin) getHelpResponse() *model.CommandResponse {
	helpText := `### KV Store Management Commands

**Available Commands:**
- ` + "`/kv set <key> <value>`" + ` - Set a key-value pair
- ` + "`/kv get <key>`" + ` - Get the value for a key
- ` + "`/kv get <pluginid:key>`" + ` - Get the value for a key from another plugin
- ` + "`/kv delete <key>`" + ` - Delete a key-value pair
- ` + "`/kv list`" + ` - List all keys in this plugin's store
- ` + "`/kv list --all`" + ` - List all keys from all plugins
- ` + "`/kv deleteall`" + ` - Delete all key-value pairs
- ` + "`/kv help`" + ` - Show this help message

**Examples:**
- ` + "`/kv set mykey Hello World`" + `
- ` + "`/kv get mykey`" + `
- ` + "`/kv get com.manybugs.mattermost-plugin-feed:some-key`" + `
- ` + "`/kv list --all`" + `
- ` + "`/kv delete mykey`" + `
`

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         helpText,
	}
}

// KVEntry represents a key-value entry with plugin information
type KVEntry struct {
	PluginID string
	Key      string
	Value    []byte
}

// isPostgreSQL checks if the database driver is PostgreSQL
func (p *Plugin) isPostgreSQL() bool {
	config := p.API.GetConfig()
	if config == nil || config.SqlSettings.DriverName == nil {
		return false
	}
	driverName := *config.SqlSettings.DriverName
	return strings.Contains(strings.ToLower(driverName), "postgres")
}

// formatSQLQuery converts MySQL-style placeholders (?) to database-specific placeholders
func (p *Plugin) formatSQLQuery(query string) string {
	if !p.isPostgreSQL() {
		return query
	}

	// Convert MySQL-style ? placeholders to PostgreSQL-style $1, $2, etc.
	var result strings.Builder
	paramNum := 1
	for _, ch := range query {
		if ch == '?' {
			result.WriteString(fmt.Sprintf("$%d", paramNum))
			paramNum++
		} else {
			result.WriteRune(ch)
		}
	}
	return result.String()
}

// handleListAll lists all keys from all plugins using database access
func (p *Plugin) handleListAll() (*model.CommandResponse, *model.AppError) {
	entries, err := p.getAllPluginKVEntries()
	if err != nil {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         fmt.Sprintf("Error listing keys from all plugins: %s", err),
		}, nil
	}

	if len(entries) == 0 {
		return &model.CommandResponse{
			ResponseType: model.CommandResponseTypeEphemeral,
			Text:         "No keys found in any plugin's KV store",
		}, nil
	}

	// Group by plugin
	pluginMap := make(map[string][]string)
	for _, entry := range entries {
		pluginMap[entry.PluginID] = append(pluginMap[entry.PluginID], entry.Key)
	}

	var text strings.Builder
	text.WriteString(fmt.Sprintf("**All Plugin KV Store Keys** (found %d keys across %d plugins):\n\n", len(entries), len(pluginMap)))

	for pluginID, keys := range pluginMap {
		text.WriteString(fmt.Sprintf("**Plugin:** `%s` (%d keys)\n", pluginID, len(keys)))
		for i, key := range keys {
			text.WriteString(fmt.Sprintf("  %d. `%s`\n", i+1, key))
		}
		text.WriteString("\n")
	}

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         text.String(),
	}, nil
}

// getAllPluginKVEntries retrieves all KV entries from all plugins using database access
func (p *Plugin) getAllPluginKVEntries() ([]KVEntry, error) {
	// Get database connection
	connID, err := p.Driver.Conn(false) // false = not master, read replica is fine
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	defer p.Driver.ConnClose(connID)

	// Query the PluginKeyValueStore table
	// The table structure is: PluginId, PKey, PValue, ExpireAt
	query := "SELECT PluginId, PKey FROM PluginKeyValueStore ORDER BY PluginId, PKey"
	rowsID, err := p.Driver.ConnQuery(connID, query, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to query KV store: %w", err)
	}
	defer p.Driver.RowsClose(rowsID)

	var entries []KVEntry
	for {
		dest := make([]driver.Value, 2)
		err := p.Driver.RowsNext(rowsID, dest)
		if err != nil {
			break
		}

		pluginID, ok := dest[0].(string)
		if !ok {
			continue
		}
		key, ok := dest[1].(string)
		if !ok {
			continue
		}

		entries = append(entries, KVEntry{
			PluginID: pluginID,
			Key:      key,
		})
	}

	return entries, nil
}

// getPluginKVValue retrieves a value for a specific plugin and key using database access
func (p *Plugin) getPluginKVValue(pluginID, key string) ([]byte, error) {
	// Get database connection
	connID, err := p.Driver.Conn(false)
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection: %w", err)
	}
	defer p.Driver.ConnClose(connID)

	// Query the PluginKeyValueStore table with database-specific placeholders
	query := p.formatSQLQuery("SELECT PValue FROM PluginKeyValueStore WHERE PluginId = ? AND PKey = ?")
	args := []driver.NamedValue{
		{Ordinal: 1, Value: pluginID},
		{Ordinal: 2, Value: key},
	}

	rowsID, err := p.Driver.ConnQuery(connID, query, args)
	if err != nil {
		return nil, fmt.Errorf("failed to query KV store: %w", err)
	}
	defer p.Driver.RowsClose(rowsID)

	dest := make([]driver.Value, 1)
	err = p.Driver.RowsNext(rowsID, dest)
	if err != nil {
		return nil, fmt.Errorf("key not found")
	}

	value, ok := dest[0].([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid value type")
	}

	return value, nil
}

// MessageWillBePosted is invoked when a message is posted by a user before it is committed
// to the database.
func (p *Plugin) MessageWillBePosted(c *plugin.Context, post *model.Post) (*model.Post, string) {
	if post.Message == "Hello, world!" {
		post.Message = post.Message + " (This message was modified by the hello-world plugin)"
	}
	return post, ""
}

func main() {
	plugin.ClientMain(&Plugin{})
}
