package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfig_GetPromptDir(t *testing.T) {
	type fields struct {
		Endpoint  string
		Model     string
		Key       string
		PromptDir string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "",
			fields: fields{
				Endpoint:  "",
				Model:     "",
				Key:       "",
				PromptDir: "",
			},
			want: filepath.Join(os.Getenv("HOME"), ".aicc", "prompt"),
		},
		{
			name: "",
			fields: fields{
				Endpoint:  "",
				Model:     "",
				Key:       "",
				PromptDir: "/opt/aicc-prompt",
			},
			want: "/opt/aicc-prompt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			receiver := &Config{
				Endpoint:  tt.fields.Endpoint,
				Model:     tt.fields.Model,
				Key:       tt.fields.Key,
				PromptDir: tt.fields.PromptDir,
			}
			if got := receiver.GetPromptDir(); got != tt.want {
				t.Errorf("GetPromptDir() = %v, want %v", got, tt.want)
			}
		})
	}
}
