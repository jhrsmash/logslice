# logslice

Fast log file slicer that filters by time range and severity without loading full files into memory.

## Installation

```bash
go install github.com/yourusername/logslice@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logslice.git && cd logslice && go build ./...
```

## Usage

```bash
# Filter logs by time range
logslice --from "2024-01-15 08:00:00" --to "2024-01-15 09:00:00" app.log

# Filter by severity level
logslice --level ERROR app.log

# Combine time range and severity
logslice --from "2024-01-15 08:00:00" --to "2024-01-15 09:00:00" --level WARN app.log

# Read from stdin
cat app.log | logslice --level ERROR
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--from` | Start of time range (RFC3339 or common formats) | — |
| `--to` | End of time range | — |
| `--level` | Minimum severity level (`DEBUG`, `INFO`, `WARN`, `ERROR`) | `DEBUG` |
| `--format` | Log timestamp format | auto-detect |

## How It Works

logslice streams log files line by line, parsing only the timestamp and severity fields needed for filtering. No full file is loaded into memory, making it suitable for very large log files.

## Requirements

- Go 1.21+

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

## License

MIT © 2024 yourusername