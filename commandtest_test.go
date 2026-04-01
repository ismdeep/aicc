package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRunConfigTests_PrintsEndpointAndModelBeforeCases(t *testing.T) {
	config := &Config{
		Endpoint: "https://example.com/api/v1",
		Model:    "test-model",
		Key:      "test-key",
	}

	var buf bytes.Buffer
	err := runConfigTests(&buf, config, func(endpoint string, model string, key string, input string) (string, error) {
		return "ok", nil
	})
	if err != nil {
		t.Fatalf("runConfigTests() error = %v", err)
	}

	output := buf.String()
	endpointIndex := strings.Index(output, "Endpoint: https://example.com/api/v1")
	modelIndex := strings.Index(output, "Model: test-model")
	testPromptIndex := strings.Index(output, "Test 1 Prompt:")

	if endpointIndex == -1 {
		t.Fatalf("output missing endpoint: %q", output)
	}
	if modelIndex == -1 {
		t.Fatalf("output missing model: %q", output)
	}
	if testPromptIndex == -1 {
		t.Fatalf("output missing test prompt: %q", output)
	}
	if endpointIndex > testPromptIndex {
		t.Fatalf("endpoint printed after test prompt: %q", output)
	}
	if modelIndex > testPromptIndex {
		t.Fatalf("model printed after test prompt: %q", output)
	}
}
