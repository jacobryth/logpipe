# logpipe

Structured log transformer that normalizes mixed-format logs from multiple services into a unified JSON stream.

## Installation

```bash
go install github.com/yourusername/logpipe@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logpipe.git && cd logpipe && go build ./...
```

## Usage

Pipe logs from any service directly into `logpipe` to normalize them into a consistent JSON format:

```bash
# Normalize logs from a single service
./my-service | logpipe

# Combine logs from multiple sources
tail -f /var/log/app.log /var/log/nginx/access.log | logpipe
```

**Example output:**

```json
{"timestamp":"2024-01-15T10:23:45Z","level":"error","service":"app","message":"connection refused","source":"/var/log/app.log"}
{"timestamp":"2024-01-15T10:23:46Z","level":"info","service":"nginx","message":"GET /health 200","source":"/var/log/nginx/access.log"}
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--pretty` | `false` | Pretty-print JSON output |
| `--service` | `""` | Override the service name field |
| `--level` | `info` | Default log level if not detected |

```bash
logpipe --pretty --service=myapp < app.log
```

## Requirements

- Go 1.21+

## License

MIT © 2024 [yourusername](https://github.com/yourusername)