package main

import (
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
			Text:         "Usage: `/kv get <key>`",
		}, nil
	}

	key := args[2]

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
- ` + "`/kv delete <key>`" + ` - Delete a key-value pair
- ` + "`/kv list`" + ` - List all keys in the store
- ` + "`/kv deleteall`" + ` - Delete all key-value pairs
- ` + "`/kv help`" + ` - Show this help message

**Examples:**
- ` + "`/kv set mykey Hello World`" + `
- ` + "`/kv get mykey`" + `
- ` + "`/kv delete mykey`" + `
`

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         helpText,
	}
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
