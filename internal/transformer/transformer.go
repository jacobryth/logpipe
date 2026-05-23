package transformer

import (
	"encoding/json"
	"fmt"
	"time"
)

// LogEntry represents a normalized log entry in unified JSON format.
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Service   string                 `json:"service,omitempty"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
}

// Transformer normalizes a parsed log record into a LogEntry.
type Transformer struct {
	service string
}

// New creates a new Transformer for the given service name.
func New(service string) *Transformer {
	return &Transformer{service: service}
}

// Transform converts a raw parsed map into a normalized LogEntry.
func (t *Transformer) Transform(record map[string]interface{}) (*LogEntry, error) {
	if record == nil {
		return nil, fmt.Errorf("transformer: nil record")
	}

	entry := &LogEntry{
		Service: t.service,
		Fields:  make(map[string]interface{}),
	}

	for k, v := range record {
		switch k {
		case "timestamp", "time", "ts", "@timestamp":
			entry.Timestamp = fmt.Sprintf("%v", v)
		case "level", "lvl", "severity":
			entry.Level = normalizeLevel(fmt.Sprintf("%v", v))
		case "message", "msg", "text":
			entry.Message = fmt.Sprintf("%v", v)
		default:
			entry.Fields[k] = v
		}
	}

	if entry.Timestamp == "" {
		entry.Timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	if entry.Level == "" {
		entry.Level = "info"
	}

	return entry, nil
}

// ToJSON serializes a LogEntry to a JSON byte slice.
func ToJSON(entry *LogEntry) ([]byte, error) {
	return json.Marshal(entry)
}

func normalizeLevel(raw string) string {
	switch raw {
	case "debug", "DEBUG", "dbg":
		return "debug"
	case "info", "INFO", "information":
		return "info"
	case "warn", "WARN", "warning", "WARNING":
		return "warn"
	case "error", "ERROR", "err", "ERR":
		return "error"
	case "fatal", "FATAL", "critical", "CRITICAL":
		return "fatal"
	default:
		return raw
	}
}
