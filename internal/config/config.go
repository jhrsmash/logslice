package config

import (
	"errors"
	"time"

	"github.com/yourorg/logslice/internal/output"
)

// Config holds all runtime configuration for a logslice run.
type Config struct {
	// FilePath is the path to the log file to slice.
	FilePath string

	// Since filters out log lines before this time (inclusive).
	// Zero value means no lower bound.
	Since time.Time

	// Until filters out log lines after this time (inclusive).
	// Zero value means no upper bound.
	Until time.Time

	// Severity is the minimum severity level to include.
	// Empty string means all severities.
	Severity string

	// Format controls output formatting.
	Format output.Format

	// ShowStats prints summary statistics to stderr after slicing.
	ShowStats bool
}

// Validate checks that the Config is internally consistent.
func (c *Config) Validate() error {
	if c.FilePath == "" {
		return errors.New("config: file path must not be empty")
	}
	if !c.Since.IsZero() && !c.Until.IsZero() && c.Until.Before(c.Since) {
		return errors.New("config: --until must not be before --since")
	}
	return nil
}
