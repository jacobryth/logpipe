package parser

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// JSONParser handles log lines that are already valid JSON objects.
type JSONParser struct{}

func (JSONParser) Name() string { return "json" }

func (JSONParser) CanParse(line string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, "{")
}

func (JSONParser) Parse(line string) (*LogEntry, error) {
	var raw map[string]any
	if err := json.Unmarshal([]byte(line), &raw); err != nil {
		return nil, fmt.Errorf("json parser: %w", err)
	}

	entry := &LogEntry{
		Timestamp: time.Now().UTC(),
		Fields:    make(map[string]any),
		Raw:       line,
	}

	known := map[string]bool{"timestamp": true, "time": true, "ts": true,
		"level": true, "severity": true, "msg": true, "message": true,
		"service": true, "logger": true}

	for k, v := range raw {
		switch k {
		case "timestamp", "time", "ts":
			if s, ok := v.(string); ok {
				if t, err := time.Parse(time.RFC3339Nano, s); err == nil {
					entry.Timestamp = t
				}
			}
		case "level", "severity":
			if s, ok := v.(string); ok {
				entry.Level = strings.ToLower(s)
			}
		case "msg", "message":
			if s, ok := v.(string); ok {
				entry.Message = s
			}
		case "service", "logger":
			if s, ok := v.(string); ok {
				entry.Service = s
			}
		}
		if !known[k] {
			entry.Fields[k] = v
		}
	}

	if entry.Message == "" {
		return nil, fmt.Errorf("json parser: missing message field")
	}
	return entry, nil
}
