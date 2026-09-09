package main

import (
	"bytes"
	"errors"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunModelsCommandWithoutCheck(t *testing.T) {
	config := &Config{
		Endpoint: "https://example.com/api/v1",
		Key:      "test-key",
	}
	models := []ModelData{
		{ID: "glm-5", OwnedBy: "modelverse"},
		{ID: "deepseek-v3.1", OwnedBy: "deepseek"},
	}

	var buf bytes.Buffer
	err := runModelsCommand(&buf, config, models, false, false, 5, func(endpoint string, model string, key string, input string) (string, error) {
		t.Fatalf("requestContent should not be called when check is disabled")
		return "", nil
	})
	if err != nil {
		t.Fatalf("runModelsCommand() error = %v", err)
	}

	output := buf.String()
	if strings.Contains(output, "STATUS") {
		t.Fatalf("output unexpectedly contains STATUS column: %q", output)
	}
	if !strings.Contains(output, "MODEL") || !strings.Contains(output, "BY") {
		t.Fatalf("output missing expected columns: %q", output)
	}
	if strings.Index(output, "deepseek-v3.1") > strings.Index(output, "glm-5") {
		t.Fatalf("output is not sorted by model name: %q", output)
	}
}

func TestRunModelsCommandWithCheck(t *testing.T) {
	config := &Config{
		Endpoint: "https://example.com/api/v1",
		Key:      "test-key",
	}
	models := []ModelData{
		{ID: "glm-5", OwnedBy: "modelverse"},
		{ID: "deepseek-v3.1", OwnedBy: "deepseek"},
	}

	var buf bytes.Buffer
	err := runModelsCommand(&buf, config, models, true, false, 5, func(endpoint string, model string, key string, input string) (string, error) {
		if endpoint != config.Endpoint {
			t.Fatalf("endpoint = %s, want %s", endpoint, config.Endpoint)
		}
		if key != config.Key {
			t.Fatalf("key = %s, want %s", key, config.Key)
		}
		if input != "Hi." {
			t.Fatalf("input = %s, want %s", input, "Hi.")
		}
		if model == "glm-5" {
			return "", errors.New("model unavailable")
		}
		return "Hi.", nil
	})
	if err != nil {
		t.Fatalf("runModelsCommand() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "STATUS") {
		t.Fatalf("output missing STATUS column: %q", output)
	}
	if !strings.Contains(output, "deepseek-v3.1") || !strings.Contains(output, "Available") {
		t.Fatalf("output missing available model status: %q", output)
	}
	if !strings.Contains(output, "glm-5") || !strings.Contains(output, "Unavailable") {
		t.Fatalf("output missing unavailable model status: %q", output)
	}
	if strings.Index(output, "deepseek-v3.1") > strings.Index(output, "glm-5") {
		t.Fatalf("output is not sorted by model name: %q", output)
	}
}

func TestRunModelsCommandWithAvailableOnly(t *testing.T) {
	config := &Config{
		Endpoint: "https://example.com/api/v1",
		Key:      "test-key",
	}
	models := []ModelData{
		{ID: "glm-5", OwnedBy: "modelverse"},
		{ID: "deepseek-v3.1", OwnedBy: "deepseek"},
	}

	var buf bytes.Buffer
	err := runModelsCommand(&buf, config, models, true, true, 5, func(endpoint string, model string, key string, input string) (string, error) {
		if model == "glm-5" {
			return "", errors.New("model unavailable")
		}
		return "Hi.", nil
	})
	if err != nil {
		t.Fatalf("runModelsCommand() error = %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "deepseek-v3.1") || !strings.Contains(output, "Available") {
		t.Fatalf("output missing available model status: %q", output)
	}
	if strings.Contains(output, "glm-5") || strings.Contains(output, "Unavailable") {
		t.Fatalf("output unexpectedly contains unavailable model: %q", output)
	}
}

func TestCommandModelsAvailableOnlyRequiresCheck(t *testing.T) {
	cmd := CommandModels()
	args := []string{"--available-only"}
	if err := cmd.ParseFlags(args); err != nil {
		t.Fatalf("ParseFlags() error = %v", err)
	}

	err := cmd.PreRunE(cmd, nil)
	if err == nil {
		t.Fatal("expected error when --available-only is used without --check")
	}
	if err.Error() != "--available-only requires --check" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCommandModelsConcurrencyMustBePositive(t *testing.T) {
	cmd := CommandModels()
	args := []string{"--concurrency", "0"}
	if err := cmd.ParseFlags(args); err != nil {
		t.Fatalf("ParseFlags() error = %v", err)
	}

	err := cmd.PreRunE(cmd, nil)
	if err == nil {
		t.Fatal("expected error when --concurrency is not positive")
	}
	if err.Error() != "--concurrency must be greater than 0" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestRunModelsCommandWithConcurrentCheck(t *testing.T) {
	config := &Config{
		Endpoint: "https://example.com/api/v1",
		Key:      "test-key",
	}
	models := []ModelData{
		{ID: "glm-5", OwnedBy: "modelverse"},
		{ID: "deepseek-v3.1", OwnedBy: "deepseek"},
		{ID: "qwen3", OwnedBy: "qwen"},
	}

	var buf bytes.Buffer
	started := make(chan struct{}, len(models))
	release := make(chan struct{})
	var current int32
	var maxConcurrent int32

	go func() {
		for range len(models) {
			<-started
		}
		close(release)
	}()

	err := runModelsCommand(&buf, config, models, true, false, len(models), func(endpoint string, model string, key string, input string) (string, error) {
		running := atomic.AddInt32(&current, 1)
		for {
			recorded := atomic.LoadInt32(&maxConcurrent)
			if running <= recorded || atomic.CompareAndSwapInt32(&maxConcurrent, recorded, running) {
				break
			}
		}

		started <- struct{}{}
		<-release
		atomic.AddInt32(&current, -1)
		return "Hi.", nil
	})
	if err != nil {
		t.Fatalf("runModelsCommand() error = %v", err)
	}

	if maxConcurrent < 2 {
		t.Fatalf("maxConcurrent = %d, want at least 2", maxConcurrent)
	}
}

func TestRunModelsCommandWritesChecksAsTheyComplete(t *testing.T) {
	config := &Config{
		Endpoint: "https://example.com/api/v1",
		Key:      "test-key",
	}
	models := []ModelData{
		{ID: "first", OwnedBy: "test"},
		{ID: "second", OwnedBy: "test"},
	}

	firstCheckDone := make(chan struct{})
	releaseSecondCheck := make(chan struct{})
	writer := &blockingWriter{
		firstRow: make(chan struct{}),
	}

	runDone := make(chan error, 1)
	go func() {
		runDone <- runModelsCommand(writer, config, models, true, false, 1, func(endpoint string, model string, key string, input string) (string, error) {
			if model == "first" {
				close(firstCheckDone)
				return "Hi.", nil
			}
			<-releaseSecondCheck
			return "Hi.", nil
		})
	}()

	select {
	case <-writer.firstRow:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for the first model row")
	}

	select {
	case <-firstCheckDone:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for the first model check")
	}

	close(releaseSecondCheck)
	if err := <-runDone; err != nil {
		t.Fatalf("runModelsCommand() error = %v", err)
	}

	output := writer.String()
	if !strings.Contains(output, "first") || !strings.Contains(output, "second") {
		t.Fatalf("output missing model rows: %q", output)
	}
}

type blockingWriter struct {
	mu       sync.Mutex
	firstRow chan struct{}
	rowSent  bool
	buf      bytes.Buffer
}

func (w *blockingWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if !w.rowSent && bytes.Contains(p, []byte("first")) {
		w.rowSent = true
		close(w.firstRow)
	}
	return w.buf.Write(p)
}

func (w *blockingWriter) String() string {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.buf.String()
}
