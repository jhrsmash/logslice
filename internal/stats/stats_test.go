package stats_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/stats"
)

func TestStats_InitialValues(t *testing.T) {
	s := stats.New()
	if s.LinesRead != 0 || s.LinesMatched != 0 || s.BytesRead != 0 {
		t.Fatal("expected all counters to start at zero")
	}
	if s.Started.IsZero() {
		t.Fatal("expected Started to be set")
	}
}

func TestStats_RecordLine(t *testing.T) {
	s := stats.New()
	s.RecordLine(42)
	s.RecordLine(8)
	if s.LinesRead != 2 {
		t.Fatalf("expected LinesRead=2, got %d", s.LinesRead)
	}
	if s.BytesRead != 50 {
		t.Fatalf("expected BytesRead=50, got %d", s.BytesRead)
	}
}

func TestStats_RecordMatch(t *testing.T) {
	s := stats.New()
	s.RecordMatch()
	s.RecordMatch()
	if s.LinesMatched != 2 {
		t.Fatalf("expected LinesMatched=2, got %d", s.LinesMatched)
	}
}

func TestStats_Finish_Idempotentent(t *testing.T) {
	s := stats.New()
	s.Finish()
	first := s.Finished
	time.Sleep(5 * time.Millisecond)
	s.Finish()
	if !s.Finished.Equal(first) {
		t.Fatal("expected Finish to be idempotent")
	}
}

func TestStats_Elapsed(t *testing.T) {
	s := stats.New()
	time.Sleep(10 * time.Millisecond)
	s.Finish()
	if s.Elapsed() < 10*time.Millisecond {
		t.Fatalf("expected elapsed >= 10ms, got %s", s.Elapsed())
	}
}

func TestStats_Write(t *testing.T) {
	s := stats.New()
	s.RecordLine(100)
	s.RecordMatch()
	s.Finish()

	var buf bytes.Buffer
	if err := s.Write(&buf); err != nil {
		t.Fatalf("Write returned error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"lines read: 1", "lines matched: 1", "bytes read: 100"} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q; got: %s", want, out)
		}
	}
}
