package config

import (
	"errors"
	"time"
)

// Config holds all runtime options for a logslice invocation.
type Config struct {
	// FilePath is the log file to slice.
	FilePath string

	// Since and Until bound the time range (both optional).
	Since *time.Time
	Until *time.Time

	// Severity filters lines to the given level and above (empty = all).
	Severity string

	// Format controls output formatting: "raw" or "json".
	Format string

	// TailN, when > 0, emits only the last N matching lines.
	TailN int

	// Follow keeps the file open and streams new lines as they arrive.
	Follow bool

	// Color enables ANSI colour highlighting.
	Color bool

	// MaxWidth truncates output lines to this rune width (0 = unlimited).
	MaxWidth int

	// RateLimit caps output at this many lines per second (0 = unlimited).
	RateLimit int
}

// Validate returns an error if the Config is logically inconsistent.
func (c *Config) Validate() error {
	if c.FilePath == "" {
		return errors.New("config: file path must not be empty")
	}
	if c.Since != nil && c.Until != nil {
		if !c.Until.After(*c.Since) {
			return errors.New("config: until must be after since")
		}
	}
	if c.RateLimit < 0 {
		return errors.New("config: rate-limit must be >= 0")
	}
	if c.MaxWidth < 0 {
		return errors.New("config: max-width must be >= 0")
	}
	return nil
}
