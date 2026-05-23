package pipeline_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/config"
	"github.com/yourorg/logslice/internal/pipeline"
)

func baseConfig() *config.Config {
	return &config.Config{
		FilePath:   "test.log",
		SampleRate: 1,
		Format:     "raw",
	}
}

func TestPipeline_NilConfig(t *testing.T) {
	_, err := pipeline.New(nil, &bytes.Buffer{})
	if err == nil {
		t.Fatal("expected error for nil config")
	}
}

func TestPipeline_NilWriter(t *testing.T) {
	_, err := pipeline.New(baseConfig(), nil)
	if err == nil {
		t.Fatal("expected error for nil writer")
	}
}

func TestPipeline_Process_ValidLine(t *testing.T) {
	var buf bytes.Buffer
	p, err := pipeline.New(baseConfig(), &buf)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	defer p.Close()

	raw := []byte("2024-01-15T10:00:00Z INFO hello world")
	n, err := p.Process(raw)
	if err != nil {
		t.Fatalf("Process: %v", err)
	}
	if n == 0 {
		t.Fatal("expected bytes written")
	}
	if !strings.Contains(buf.String(), "hello world") {
		t.Errorf("output missing message: %q", buf.String())
	}
}

func TestPipeline_Process_Unparseable(t *testing.T) {
	var buf bytes.Buffer
	p, _ := pipeline.New(baseConfig(), &buf)
	defer p.Close()

	n, err := p.Process([]byte("not a log line at all!!!"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if n != 0 {
		t.Errorf("expected 0 bytes written for unparseable line, got %d", n)
	}
}

func TestPipeline_Process_SeverityFilter(t *testing.T) {
	cfg := baseConfig()
	cfg.MinSeverity = "ERROR"

	var buf bytes.Buffer
	p, _ := pipeline.New(cfg, &buf)
	defer p.Close()

	p.Process([]byte("2024-01-15T10:00:00Z DEBUG low priority"))
	if buf.Len() != 0 {
		t.Errorf("DEBUG line should have been filtered, got %q", buf.String())
	}

	p.Process([]byte("2024-01-15T10:00:01Z ERROR something broke"))
	if buf.Len() == 0 {
		t.Error("ERROR line should have passed severity filter")
	}
}

func TestPipeline_Process_TimeRangeFilter(t *testing.T) {
	cfg := baseConfig()
	cfg.Since = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)
	cfg.Until = time.Date(2024, 1, 15, 13, 0, 0, 0, time.UTC)

	var buf bytes.Buffer
	p, _ := pipeline.New(cfg, &buf)
	defer p.Close()

	p.Process([]byte("2024-01-15T10:00:00Z INFO too early"))
	if buf.Len() != 0 {
		t.Errorf("line before Since should be filtered, got %q", buf.String())
	}

	p.Process([]byte("2024-01-15T12:30:00Z INFO within range"))
	if buf.Len() == 0 {
		t.Error("line within range should pass")
	}
}

func TestPipeline_Stats_Counts(t *testing.T) {
	var buf bytes.Buffer
	p, _ := pipeline.New(baseConfig(), &buf)
	defer p.Close()

	p.Process([]byte("not parseable"))
	p.Process([]byte("2024-01-15T10:00:00Z INFO first"))
	p.Process([]byte("2024-01-15T10:00:01Z INFO second"))

	s := p.Stats()
	if s.Lines() != 2 {
		t.Errorf("Lines: want 2, got %d", s.Lines())
	}
	if s.Matches() != 2 {
		t.Errorf("Matches: want 2, got %d", s.Matches())
	}
}
