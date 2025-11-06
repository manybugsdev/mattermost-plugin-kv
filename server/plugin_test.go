package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateKey(t *testing.T) {
	tests := []struct {
		name    string
		key     string
		wantErr bool
	}{
		{
			name:    "valid key",
			key:     "mykey",
			wantErr: false,
		},
		{
			name:    "valid key with numbers",
			key:     "mykey123",
			wantErr: false,
		},
		{
			name:    "valid key with underscore",
			key:     "my_key",
			wantErr: false,
		},
		{
			name:    "valid key with dash",
			key:     "my-key",
			wantErr: false,
		},
		{
			name:    "empty key",
			key:     "",
			wantErr: true,
		},
		{
			name:    "key with slash",
			key:     "my/key",
			wantErr: true,
		},
		{
			name:    "key with backslash",
			key:     "my\\key",
			wantErr: true,
		},
		{
			name:    "key with double dots",
			key:     "my..key",
			wantErr: true,
		},
		{
			name:    "key too long",
			key:     string(make([]byte, 257)),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateKey(tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateKey() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPluginManifest(t *testing.T) {
	p := &Plugin{}
	m := p.Manifest()

	require.NotNil(t, m, "Manifest should not be nil")
	assert.Equal(t, "com.manybugs.mattermost-plugin-kv", m.Id, "Plugin ID should match")
	assert.Equal(t, "KV Manager", m.Name, "Plugin name should match")
	assert.NotEmpty(t, m.Version, "Plugin version should not be empty")
}
