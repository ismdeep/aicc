package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestModels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Fatalf("method = %s, want %s", r.Method, http.MethodGet)
		}

		if r.URL.Path != "/api/v1/models" {
			t.Fatalf("path = %s, want %s", r.URL.Path, "/api/v1/models")
		}

		if got := r.Header.Get("Authorization"); got != "Bearer test-key" {
			t.Fatalf("authorization = %s, want %s", got, "Bearer test-key")
		}

		_, _ = io.WriteString(w, `{"data":[{"created":1761730258,"id":"deepseek-v3.1","object":"model","owned by":"deepseek"},{"created":1761730264,"id":"glm-5","object":"model","owned by":"modelverse"}]}`)
	}))
	defer server.Close()

	models, err := RequestModels(server.URL+"/api/v1", "test-key")
	if err != nil {
		t.Fatalf("RequestModels() error = %v", err)
	}

	if len(models) != 2 {
		t.Fatalf("len(models) = %d, want 2", len(models))
	}

	if models[0].ID != "deepseek-v3.1" {
		t.Fatalf("models[0].ID = %s, want %s", models[0].ID, "deepseek-v3.1")
	}

	if models[1].OwnedBy != "modelverse" {
		t.Fatalf("models[1].OwnedBy = %s, want %s", models[1].OwnedBy, "modelverse")
	}
}

func TestJoinURL(t *testing.T) {
	tests := []struct {
		name     string
		endpoint string
		path     string
		want     string
	}{
		{
			name:     "without trailing slash",
			endpoint: "https://example.com/api/v1",
			path:     "/models",
			want:     "https://example.com/api/v1/models",
		},
		{
			name:     "with trailing slash",
			endpoint: "https://example.com/api/v1/",
			path:     "/models",
			want:     "https://example.com/api/v1/models",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := joinURL(tt.endpoint, tt.path); got != tt.want {
				t.Fatalf("joinURL() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestWriteModelsTable(t *testing.T) {
	var buf bytes.Buffer
	models := []ModelData{
		{ID: "deepseek-v3.1", OwnedBy: "deepseek"},
		{ID: "glm-5", OwnedBy: "modelverse"},
	}

	if err := writeModelsTable(&buf, models, nil); err != nil {
		t.Fatalf("writeModelsTable() error = %v", err)
	}

	got := buf.String()
	lines := strings.Split(strings.TrimSpace(got), "\n")
	if len(lines) != 3 {
		t.Fatalf("line count = %d, want 3, output = %q", len(lines), got)
	}

	if !strings.Contains(lines[0], "MODEL") || !strings.Contains(lines[0], "BY") {
		t.Fatalf("header = %q, want columns MODEL and BY", lines[0])
	}

	if !strings.Contains(lines[1], "deepseek-v3.1") || !strings.Contains(lines[1], "deepseek") {
		t.Fatalf("row 1 = %q, want model data", lines[1])
	}

	if !strings.Contains(lines[2], "glm-5") || !strings.Contains(lines[2], "modelverse") {
		t.Fatalf("row 2 = %q, want model data", lines[2])
	}
}
