package parser_test

import (
	"testing"
	"time"

	"github.com/yourorg/logpipe/internal/parser"
)

func newRegistry() *parser.Registry {
	r := parser.NewRegistry()
	r.Register(parser.JSONParser{})
	r.Register(parser.LogfmtParser{})
	return r
}

func TestJSONParser(t *testing.T) {
	r := newRegistry()
	line := `{"timestamp":"2024-01-15T10:00:00Z","level":"info","message":"server started","service":"api","port":8080}`

	entry, err := r.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Level != "info" {
		t.Errorf("level: got %q, want %q", entry.Level, "info")
	}
	if entry.Message != "server started" {
		t.Errorf("message: got %q, want %q", entry.Message, "server started")
	}
	if entry.Service != "api" {
		t.Errorf("service: got %q, want %q", entry.Service, "api")
	}
	expected := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	if !entry.Timestamp.Equal(expected) {
		t.Errorf("timestamp: got %v, want %v", entry.Timestamp, expected)
	}
	if _, ok := entry.Fields["port"]; !ok {
		t.Error("expected 'port' in fields")
	}
}

func TestLogfmtParser(t *testing.T) {
	r := newRegistry()
	line := `ts=2024-01-15T10:00:00Z level=warn msg="disk usage high" service=monitor usage=92`

	entry, err := r.Parse(line)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Level != "warn" {
		t.Errorf("level: got %q, want %q", entry.Level, "warn")
	}
	if entry.Message != "disk usage high" {
		t.Errorf("message: got %q, want %q", entry.Message, "disk usage high")
	}
	if entry.Service != "monitor" {
		t.Errorf("service: got %q, want %q", entry.Service, "monitor")
	}
	if _, ok := entry.Fields["usage"]; !ok {
		t.Error("expected 'usage' in fields")
	}
}

func TestFallbackParser(t *testing.T) {
	r := newRegistry()
	line := "plain text log line with no structure"

	entry, err := r.Parse(line)
	if err == nil {
		t.Fatal("expected fallback error, got nil")
	}
	if entry == nil {
		t.Fatal("expected fallback entry, got nil")
	}
	if entry.Message != line {
		t.Errorf("fallback message: got %q, want %q", entry.Message, line)
	}
	if entry.Level != "unknown" {
		t.Errorf("fallback level: got %q, want %q", entry.Level, "unknown")
	}
}
