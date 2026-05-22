package config_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/config"
	"github.com/yourorg/logslice/internal/output"
)

func baseConfig() *config.Config {
	return &config.Config{
		FilePath: "/tmp/app.log",
		Format:   output.FormatRaw,
	}
}

func TestConfig_Valid(t *testing.T) {
	cfg := baseConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestConfig_EmptyFilePath(t *testing.T) {
	cfg := baseConfig()
	cfg.FilePath = ""
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error for empty file path")
	}
}

func TestConfig_UntilBeforeSince(t *testing.T) {
	cfg := baseConfig()
	now := time.Now()
	cfg.Since = now
	cfg.Until = now.Add(-time.Minute)
	if err := cfg.Validate(); err == nil {
		t.Fatal("expected error when until is before since")
	}
}

func TestConfig_EqualSinceUntil(t *testing.T) {
	cfg := baseConfig()
	now := time.Now()
	cfg.Since = now
	cfg.Until = now
	if err := cfg.Validate(); err != nil {
		t.Fatalf("equal since/until should be valid, got: %v", err)
	}
}

func TestConfig_OnlySince(t *testing.T) {
	cfg := baseConfig()
	cfg.Since = time.Now()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("only since set should be valid, got: %v", err)
	}
}

func TestConfig_ShowStats(t *testing.T) {
	cfg := baseConfig()
	cfg.ShowStats = true
	if err := cfg.Validate(); err != nil {
		t.Fatalf("show stats flag should not affect validation, got: %v", err)
	}
}
