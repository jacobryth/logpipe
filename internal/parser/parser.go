package parser

import (
	"fmt"
	"time"
)

// LogEntry represents a normalized log record.
type LogEntry struct {
	Timestamp time.Time         `json:"timestamp"`
	Level     string            `json:"level"`
	Message   string            `json:"message"`
	Service   string            `json:"service,omitempty"`
	Fields    map[string]any    `json:"fields,omitempty"`
	Raw       string            `json:"raw,omitempty"`
}

// Parser is implemented by any format-specific log parser.
type Parser interface {
	// Name returns a human-readable identifier for the parser.
	Name() string
	// CanParse reports whether the parser recognises the given raw line.
	CanParse(line string) bool
	// Parse converts a raw log line into a LogEntry.
	Parse(line string) (*LogEntry, error)
}

// Registry holds a prioritised list of parsers.
type Registry struct {
	parsers []Parser
}

// NewRegistry returns an empty Registry.
func NewRegistry() *Registry {
	return &Registry{}
}

// Register appends a parser to the registry.
func (r *Registry) Register(p Parser) {
	r.parsers = append(r.parsers, p)
}

// Parse iterates registered parsers and returns the first successful result.
// If no parser matches, a fallback entry carrying the raw line is returned.
func (r *Registry) Parse(line string) (*LogEntry, error) {
	for _, p := range r.parsers {
		if p.CanParse(line) {
			return p.Parse(line)
		}
	}
	return &LogEntry{
		Timestamp: time.Now().UTC(),
		Level:     "unknown",
		Message:   line,
		Raw:       line,
	}, fmt.Errorf("no parser matched line")
}
