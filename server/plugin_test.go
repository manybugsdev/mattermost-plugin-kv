package main

import (
	"testing"

	"github.com/mattermost/mattermost/server/public/model"
	"github.com/mattermost/mattermost/server/public/plugin"
	"github.com/mattermost/mattermost/server/public/plugin/plugintest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOnActivate(t *testing.T) {
	p := &Plugin{}
	api := &plugintest.API{}

	// Expect the LogInfo call
	api.On("LogInfo", "Hello World plugin has been activated!").Return(nil)

	p.SetAPI(api)

	err := p.OnActivate()
	assert.NoError(t, err)

	api.AssertExpectations(t)
}

func TestMessageHasBeenPosted(t *testing.T) {
	p := &Plugin{}
	api := &plugintest.API{}

	post := &model.Post{
		UserId:    "user123",
		ChannelId: "channel123",
		Message:   "Test message",
	}

	// Expect the LogInfo call
	api.On("LogInfo", "User user123 posted: Test message").Return(nil)

	// Expect the SendEphemeralPost call
	api.On("SendEphemeralPost", "user123", mock.MatchedBy(func(p *model.Post) bool {
		return p.UserId == "user123" &&
			p.ChannelId == "channel123" &&
			p.Message == "Hello from Hello World plugin!"
	})).Return(&model.Post{})

	p.SetAPI(api)

	p.MessageHasBeenPosted(&plugin.Context{}, post)

	api.AssertExpectations(t)
}
