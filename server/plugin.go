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
	if err := p.API.RegisterCommand(&model.Command{
		Trigger:          "kv",
		DisplayName:      "KV Manager",
		Description:      "Manage Mattermost KV store",
		AutoComplete:     true,
		AutoCompleteDesc: "Available commands: set, get, delete, list, help",
		AutoCompleteHint: "[set|get|delete|list|help]",
	}); err != nil {
		return err
	}
	return nil
}

// ExecuteCommand handles slash command execution
func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	trigger := strings.TrimPrefix(strings.Fields(args.Command)[0], "/")

	if trigger != "kv" {
		return &model.CommandResponse{}, nil
	}

	// Parse command arguments
	parts := strings.Fields(args.Command)
	if len(parts) < 2 {
		return p.sendHelpResponse(), nil
	}

	subcommand := parts[1]

	switch subcommand {
	case "set":
		return p.handleSet(parts[2:])
	case "get":
		return p.handleGet(parts[2:])
	case "delete":
		return p.handleDelete(parts[2:])
	case "list":
		return p.handleList(parts[2:])
	case "help":
		return p.sendHelpResponse(), nil
	default:
		return p.sendHelpResponse(), nil
	}
}

// handleSet handles the 'set' subcommand
func (p *Plugin) handleSet(args []string) (*model.CommandResponse, *model.AppError) {
	if len(args) < 2 {
		return p.sendErrorResponse("Usage: /kv set <key> <value>"), nil
	}

	key := args[0]
	value := strings.Join(args[1:], " ")

	if err := p.API.KVSet(key, []byte(value)); err != nil {
		return p.sendErrorResponse(fmt.Sprintf("Error setting key: %v", err)), nil
	}

	return p.sendSuccessResponse(fmt.Sprintf("Successfully set key '%s' to value '%s'", key, value)), nil
}

// handleGet handles the 'get' subcommand
func (p *Plugin) handleGet(args []string) (*model.CommandResponse, *model.AppError) {
	if len(args) < 1 {
		return p.sendErrorResponse("Usage: /kv get <key>"), nil
	}

	key := args[0]

	value, err := p.API.KVGet(key)
	if err != nil {
		return p.sendErrorResponse(fmt.Sprintf("Error getting key: %v", err)), nil
	}

	if value == nil {
		return p.sendErrorResponse(fmt.Sprintf("Key '%s' not found", key)), nil
	}

	return p.sendSuccessResponse(fmt.Sprintf("Value for key '%s': %s", key, string(value))), nil
}

// handleDelete handles the 'delete' subcommand
func (p *Plugin) handleDelete(args []string) (*model.CommandResponse, *model.AppError) {
	if len(args) < 1 {
		return p.sendErrorResponse("Usage: /kv delete <key>"), nil
	}

	key := args[0]

	if err := p.API.KVDelete(key); err != nil {
		return p.sendErrorResponse(fmt.Sprintf("Error deleting key: %v", err)), nil
	}

	return p.sendSuccessResponse(fmt.Sprintf("Successfully deleted key '%s'", key)), nil
}

// handleList handles the 'list' subcommand
func (p *Plugin) handleList(args []string) (*model.CommandResponse, *model.AppError) {
	var prefix string
	if len(args) > 0 {
		prefix = args[0]
	}

	keys, err := p.API.KVList(0, 100)
	if err != nil {
		return p.sendErrorResponse(fmt.Sprintf("Error listing keys: %v", err)), nil
	}

	if len(keys) == 0 {
		return p.sendSuccessResponse("No keys found in the KV store"), nil
	}

	// Filter keys by prefix if provided
	filteredKeys := []string{}
	for _, key := range keys {
		if prefix == "" || strings.HasPrefix(key, prefix) {
			filteredKeys = append(filteredKeys, key)
		}
	}

	if len(filteredKeys) == 0 {
		return p.sendSuccessResponse(fmt.Sprintf("No keys found with prefix '%s'", prefix)), nil
	}

	response := fmt.Sprintf("Found %d key(s):\n", len(filteredKeys))
	for _, key := range filteredKeys {
		response += fmt.Sprintf("- %s\n", key)
	}

	return p.sendSuccessResponse(response), nil
}

// sendHelpResponse returns a help message
func (p *Plugin) sendHelpResponse() *model.CommandResponse {
	helpText := `### KV Manager Plugin Commands

Available commands:
- **/kv set <key> <value>** - Set a key-value pair
- **/kv get <key>** - Get the value of a key
- **/kv delete <key>** - Delete a key
- **/kv list [prefix]** - List all keys (optionally filtered by prefix)
- **/kv help** - Display this help message

Examples:
- /kv set mykey myvalue
- /kv get mykey
- /kv delete mykey
- /kv list
- /kv list my`

	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         helpText,
	}
}

// sendSuccessResponse returns a success message
func (p *Plugin) sendSuccessResponse(message string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         "✅ " + message,
	}
}

// sendErrorResponse returns an error message
func (p *Plugin) sendErrorResponse(message string) *model.CommandResponse {
	return &model.CommandResponse{
		ResponseType: model.CommandResponseTypeEphemeral,
		Text:         "❌ " + message,
	}
}

func main() {
	plugin.ClientMain(&Plugin{})
}
