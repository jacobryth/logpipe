package parser

import (
	"fmt"
	"strings"
	"time"

	"github.com/kr/logfmt"
)

// LogfmtParser handles key=value style log lines (logfmt).
type LogfmtParser struct{}

func (LogfmtParser) Name() string { return "logfmt" }

func (LogfmtParser) CanParse(line string) bool {
	// A logfmt line typically contains at least one key=value pair.
	return strings.Contains(line, "=") && !strings.HasPrefix(strings.TrimSpace(line), "{")
}

type logfmtHandler struct {
	fields map[string]string
}

func (h *logfmtHandler) HandleLogfmt(key, val []byte) error {
	h.fields[string(key)] = string(val)
	return nil
}

func (LogfmtParser) Parse(line string) (*LogEntry, error) {
	h := &logfmtHandler{fields: make(map[string]string)}
	if err := logfmt.Unmarshal([]byte(line), h); err != nil {
		return nil, fmt.Errorf("logfmt parser: %w", err)
	}

	entry := &LogEntry{
		Timestamp: time.Now().UTC(),
		Fields:    make(map[string]any),
		Raw:       line,
	}

	known := map[string]bool{"ts": true, "time": true, "level": true,
		"msg": true, "message": true, "service": true}

	for k, v := range h.fields {
		switch k {
		case "ts", "time":
			if t, err := time.Parse(time.RFC3339Nano, v); err == nil {
				entry.Timestamp = t
			}
		case "level":
			entry.Level = strings.ToLower(v)
		case "msg", "message":
			entry.Message = v
		case "service":
			entry.Service = v
		}
		if !known[k] {
			entry.Fields[k] = v
		}
	}

	if entry.Message == "" {
		return nil, fmt.Errorf("logfmt parser: missing msg/message field")
	}
	return entry, nil
}
