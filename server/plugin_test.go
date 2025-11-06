package main

import (
	"testing"
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
