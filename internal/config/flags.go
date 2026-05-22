package config

import (
	"errors"
	"flag"
	"time"
)

const timeLayout = time.RFC3339

// RegisterFlags adds logslice's CLI flags to fs.
func RegisterFlags(fs *flag.FlagSet) {
	fs.String("since", "", "include lines at or after this RFC3339 timestamp")
	fs.String("until", "", "include lines before this RFC3339 timestamp")
	fs.String("severity", "", "minimum severity level (DEBUG|INFO|WARN|ERROR)")
	fs.String("format", "raw", "output format: raw or json")
	fs.Int("tail", 0, "emit only the last N matching lines")
	fs.Bool("follow", false, "keep file open and stream new lines")
	fs.Bool("color", false, "enable ANSI colour highlighting")
	fs.Int("max-width", 0, "truncate output lines to this rune width (0 = off)")
	fs.Int("rate-limit", 0, "max output lines per second (0 = unlimited)")
}

// FromFlags builds a Config from a parsed FlagSet and positional args.
// args[0] must be the log file path.
func FromFlags(fs *flag.FlagSet, args []string) (*Config, error) {
	if len(args) == 0 {
		return nil, errors.New("flags: file path argument is required")
	}

	cfg := &Config{
		FilePath:  args[0],
		Format:    fs.Lookup("format").Value.String(),
		Severity:  fs.Lookup("severity").Value.String(),
		Follow:    fs.Lookup("follow").Value.String() == "true",
		Color:     fs.Lookup("color").Value.String() == "true",
	}

	if v := fs.Lookup("tail").Value.String(); v != "0" {
		var n int
		if _, err := parseIntFlag(v, &n); err != nil {
			return nil, err
		}
		cfg.TailN = n
	}

	if v := fs.Lookup("max-width").Value.String(); v != "0" {
		var n int
		if _, err := parseIntFlag(v, &n); err != nil {
			return nil, err
		}
		cfg.MaxWidth = n
	}

	if v := fs.Lookup("rate-limit").Value.String(); v != "0" {
		var n int
		if _, err := parseIntFlag(v, &n); err != nil {
			return nil, err
		}
		cfg.RateLimit = n
	}

	if s := fs.Lookup("since").Value.String(); s != "" {
		t, err := time.Parse(timeLayout, s)
		if err != nil {
			return nil, errors.New("flags: invalid --since: " + err.Error())
		}
		cfg.Since = &t
	}

	if s := fs.Lookup("until").Value.String(); s != "" {
		t, err := time.Parse(timeLayout, s)
		if err != nil {
			return nil, errors.New("flags: invalid --until: " + err.Error())
		}
		cfg.Until = &t
	}

	return cfg, nil
}

func parseIntFlag(s string, dst *int) (int, error) {
	var n int
	if _, err := fmt.Sscanf(s, "%d", &n); err != nil {
		return 0, errors.New("flags: expected integer, got: " + s)
	}
	*dst = n
	return n, nil
}
