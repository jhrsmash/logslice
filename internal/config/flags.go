package config

import (
	"flag"
	"fmt"
	"time"

	"github.com/yourorg/logslice/internal/output"
)

const timeLayout = time.RFC3339

// FromFlags parses a Config from the provided FlagSet after it has been parsed.
// Call fs.Parse(args) before calling FromFlags.
func FromFlags(fs *flag.FlagSet) (*Config, error) {
	cfg := &Config{}

	sinceStr := fs.Lookup("since")
	untilStr := fs.Lookup("until")
	formatStr := fs.Lookup("format")

	if sinceStr != nil && sinceStr.Value.String() != "" {
		t, err := time.Parse(timeLayout, sinceStr.Value.String())
		if err != nil {
			return nil, fmt.Errorf("config: invalid --since: %w", err)
		}
		cfg.Since = t
	}

	if untilStr != nil && untilStr.Value.String() != "" {
		t, err := time.Parse(timeLayout, untilStr.Value.String())
		if err != nil {
			return nil, fmt.Errorf("config: invalid --until: %w", err)
		}
		cfg.Until = t
	}

	if formatStr != nil {
		f, err := output.ParseFormat(formatStr.Value.String())
		if err != nil {
			return nil, fmt.Errorf("config: invalid --format: %w", err)
		}
		cfg.Format = f
	}

	if sev := fs.Lookup("severity"); sev != nil {
		cfg.Severity = sev.Value.String()
	}

	if stats := fs.Lookup("stats"); stats != nil {
		cfg.ShowStats = stats.Value.String() == "true"
	}

	if args := fs.Args(); len(args) > 0 {
		cfg.FilePath = args[0]
	}

	return cfg, nil
}

// RegisterFlags registers all logslice flags onto the given FlagSet.
func RegisterFlags(fs *flag.FlagSet) {
	fs.String("since", "", "include lines at or after this RFC3339 timestamp")
	fs.String("until", "", "include lines at or before this RFC3339 timestamp")
	fs.String("severity", "", "minimum severity level (DEBUG, INFO, WARN, ERROR)")
	fs.String("format", "raw", "output format: raw or json")
	fs.Bool("stats", false, "print summary statistics to stderr")
}
