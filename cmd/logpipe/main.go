// Package main is the entry point for the logpipe CLI tool.
// It reads log lines from stdin, parses them using the appropriate
// parser, transforms them into a unified format, and writes JSON to stdout.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/yourorg/logpipe/internal/parser"
	"github.com/yourorg/logpipe/internal/transformer"
)

const version = "0.1.0"

func main() {
	var (
		showVersion = flag.Bool("version", false, "print version and exit")
		service     = flag.String("service", "", "override service name in output records")
		pretty      = flag.Bool("pretty", false, "pretty-print JSON output")
	)
	flag.Parse()

	if *showVersion {
		fmt.Printf("logpipe %s\n", version)
		os.Exit(0)
	}

	// Build the parser registry with all supported formats.
	reg := parser.NewRegistry()

	// Build the transformer, optionally overriding the service name.
	opts := []transformer.Option{}
	if *service != "" {
		opts = append(opts, transformer.WithService(*service))
	}
	if *pretty {
		opts = append(opts, transformer.WithPretty())
	}
	t := transformer.New(opts...)

	// Determine input source: file args or stdin.
	var readers []io.Reader
	if args := flag.Args(); len(args) > 0 {
		for _, path := range args {
			f, err := os.Open(path)
			if err != nil {
				log.Fatalf("logpipe: cannot open %q: %v", path, err)
			}
			defer f.Close() //nolint:gocritic
			readers = append(readers, f)
		}
	} else {
		readers = append(readers, os.Stdin)
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	for _, r := range readers {
		if err := processReader(r, reg, t, w); err != nil {
			log.Fatalf("logpipe: %v", err)
		}
	}
}

// processReader reads lines from r, parses each one, transforms it, and
// writes the resulting JSON line to w.
func processReader(r io.Reader, reg *parser.Registry, t *transformer.Transformer, w io.Writer) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024) // up to 1 MiB per line

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}

		record, err := reg.Parse(line)
		if err != nil {
			// Unparseable lines are emitted as plain-message records so
			// downstream consumers always receive valid JSON.
			record = parser.FallbackRecord(line)
		}

		out, err := t.ToJSON(record)
		if err != nil {
			return fmt.Errorf("transform: %w", err)
		}

		if _, err := fmt.Fprintln(w, string(out)); err != nil {
			return fmt.Errorf("write: %w", err)
		}
	}

	return scanner.Err()
}
