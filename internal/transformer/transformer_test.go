package transformer

import (
	"encoding/json"
	"testing"
)

func TestTransformFullRecord(t *testing.T) {
	tr := New("auth-service")
	record := map[string]interface{}{
		"ts":      "2024-01-15T10:00:00Z",
		"lvl":     "ERROR",
		"msg":     "connection refused",
		"host":    "db-01",
		"attempt": 3,
	}

	entry, err := tr.Transform(record)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Timestamp != "2024-01-15T10:00:00Z" {
		t.Errorf("expected timestamp '2024-01-15T10:00:00Z', got %q", entry.Timestamp)
	}
	if entry.Level != "error" {
		t.Errorf("expected level 'error', got %q", entry.Level)
	}
	if entry.Message != "connection refused" {
		t.Errorf("expected message 'connection refused', got %q", entry.Message)
	}
	if entry.Service != "auth-service" {
		t.Errorf("expected service 'auth-service', got %q", entry.Service)
	}
	if entry.Fields["host"] != "db-01" {
		t.Errorf("expected fields.host 'db-01', got %v", entry.Fields["host"])
	}
}

func TestTransformDefaults(t *testing.T) {
	tr := New("")
	record := map[string]interface{}{
		"message": "hello world",
	}

	entry, err := tr.Transform(record)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if entry.Level != "info" {
		t.Errorf("expected default level 'info', got %q", entry.Level)
	}
	if entry.Timestamp == "" {
		t.Error("expected non-empty default timestamp")
	}
}

func TestTransformNilRecord(t *testing.T) {
	tr := New("svc")
	_, err := tr.Transform(nil)
	if err == nil {
		t.Error("expected error for nil record")
	}
}

func TestToJSON(t *testing.T) {
	entry := &LogEntry{
		Timestamp: "2024-01-15T10:00:00Z",
		Level:     "warn",
		Message:   "disk usage high",
		Service:   "monitor",
		Fields:    map[string]interface{}{"usage": 92},
	}

	b, err := ToJSON(entry)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var out map[string]interface{}
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatalf("invalid JSON output: %v", err)
	}
	if out["level"] != "warn" {
		t.Errorf("expected level 'warn' in JSON, got %v", out["level"])
	}
}

func TestNormalizeLevel(t *testing.T) {
	cases := map[string]string{
		"DEBUG": "debug", "dbg": "debug",
		"INFO": "info", "information": "info",
		"WARN": "warn", "warning": "warn",
		"ERR": "error", "ERROR": "error",
		"FATAL": "fatal", "critical": "fatal",
		"custom": "custom",
	}
	for input, want := range cases {
		if got := normalizeLevel(input); got != want {
			t.Errorf("normalizeLevel(%q) = %q, want %q", input, got, want)
		}
	}
}
