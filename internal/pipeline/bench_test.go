package pipeline_test

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/pipeline"
)

// BenchmarkPipeline_Process measures throughput of the full processing chain
// for lines that pass all filters (worst-case path: every stage executes).
func BenchmarkPipeline_Process(b *testing.B) {
	var buf bytes.Buffer
	p, err := pipeline.New(baseConfig(), &buf)
	if err != nil {
		b.Fatalf("New: %v", err)
	}
	defer p.Close()

	raw := []byte("2024-01-15T10:00:00Z INFO benchmark log message payload")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		p.Process(raw)
	}
}

// BenchmarkPipeline_Process_Filtered measures throughput when lines are
// rejected early by the severity filter.
func BenchmarkPipeline_Process_Filtered(b *testing.B) {
	cfg := baseConfig()
	cfg.MinSeverity = "ERROR"

	var buf bytes.Buffer
	p, _ := pipeline.New(cfg, &buf)
	defer p.Close()

	raw := []byte("2024-01-15T10:00:00Z DEBUG this will be filtered")

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		p.Process(raw)
	}
}

// BenchmarkPipeline_Process_HighVolume exercises the pipeline with many
// distinct timestamps to avoid dedupe short-circuiting.
func BenchmarkPipeline_Process_HighVolume(b *testing.B) {
	var buf bytes.Buffer
	p, _ := pipeline.New(baseConfig(), &buf)
	defer p.Close()

	base := time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	lines := make([][]byte, 1000)
	for i := range lines {
		ts := base.Add(time.Duration(i) * time.Second)
		lines[i] = []byte(fmt.Sprintf("%s INFO message number %d", ts.Format(time.RFC3339), i))
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		p.Process(lines[i%len(lines)])
	}
}
