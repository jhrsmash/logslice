package config_test

import (
	"flag"
	"testing"

	"github.com/yourorg/logslice/internal/config"
	"github.com/yourorg/logslice/internal/output"
)

func makeFS(args []string) (*flag.FlagSet, error) {
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	config.RegisterFlags(fs)
	return fs, fs.Parse(args)
}

func TestFromFlags_Defaults(t *testing.T) {
	fs, err := makeFS([]string{"/var/log/app.log"})
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := config.FromFlags(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.FilePath != "/var/log/app.log" {
		t.Errorf("expected file path, got %q", cfg.FilePath)
	}
	if cfg.Format != output.FormatRaw {
		t.Errorf("expected raw format, got %v", cfg.Format)
	}
	if cfg.ShowStats {
		t.Error("expected ShowStats=false by default")
	}
}

func TestFromFlags_Since(t *testing.T) {
	fs, err := makeFS([]string{"--since", "2024-01-15T10:00:00Z", "app.log"})
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := config.FromFlags(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Since.IsZero() {
		t.Error("expected Since to be set")
	}
}

func TestFromFlags_InvalidSince(t *testing.T) {
	fs, err := makeFS([]string{"--since", "not-a-time", "app.log"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = config.FromFlags(fs)
	if err == nil {
		t.Fatal("expected error for invalid --since")
	}
}

func TestFromFlags_JSONFormat(t *testing.T) {
	fs, err := makeFS([]string{"--format", "json", "app.log"})
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := config.FromFlags(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Format != output.FormatJSON {
		t.Errorf("expected json format, got %v", cfg.Format)
	}
}

func TestFromFlags_Stats(t *testing.T) {
	fs, err := makeFS([]string{"--stats", "app.log"})
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := config.FromFlags(fs)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !cfg.ShowStats {
		t.Error("expected ShowStats=true")
	}
}
