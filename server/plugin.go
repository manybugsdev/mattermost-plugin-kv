package main

import (
	"fmt"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin
}

// OnActivate is invoked when the plugin is activated.
func (p *Plugin) OnActivate() error {
	p.API.LogInfo("Hello World plugin has been activated!")
	return nil
}

// MessageHasBeenPosted is called after a message has been posted to the database.
// This hook demonstrates how a plugin can react to user messages.
func (p *Plugin) MessageHasBeenPosted(c *plugin.Context, post *model.Post) {
	// Log information about the posted message
	p.API.LogInfo(fmt.Sprintf("User %s posted: %s", post.UserId, post.Message))

	// Send an ephemeral message (visible only to the user who posted) as a response
	ephemeralPost := &model.Post{
		UserId:    post.UserId,
		ChannelId: post.ChannelId,
		Message:   "Hello from Hello World plugin!",
	}
	p.API.SendEphemeralPost(post.UserId, ephemeralPost)
}

// See https://developers.mattermost.com/extend/plugins/server/reference/
