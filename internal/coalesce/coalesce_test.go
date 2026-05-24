package coalesce_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/coalesce"
	"github.com/yourorg/logslice/internal/parser"
)

var base = time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)

func line(msg string, ts time.Time) *parser.LogLine {
	return &parser.LogLine{Timestamp: ts, Severity: "INFO", Message: msg}
}

func contLine(msg string) *parser.LogLine {
	// Continuation lines carry a zero timestamp.
	return &parser.LogLine{Message: msg}
}

func TestCoalesce_NilLine_Ignored(t *testing.T) {
	c := coalesce.New(0)
	if got := c.Push(nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestCoalesce_SingleLine_FlushReturnsIt(t *testing.T) {
	c := coalesce.New(0)
	if flushed := c.Push(line("hello", base)); flushed != nil {
		t.Fatalf("first push should not flush, got %v", flushed)
	}
	got := c.Flush()
	if got == nil || got.Message != "hello" {
		t.Fatalf("expected 'hello', got %v", got)
	}
}

func TestCoalesce_ContinuationAppended(t *testing.T) {
	c := coalesce.New(0)
	c.Push(line("first", base))
	c.Push(contLine("  at foo.go:10"))
	c.Push(contLine("  at bar.go:20"))

	got := c.Flush()
	if got == nil {
		t.Fatal("expected merged line")
	}
	want := "first\n  at foo.go:10\n  at bar.go:20"
	if got.Message != want {
		t.Fatalf("message mismatch\ngot:  %q\nwant: %q", got.Message, want)
	}
}

func TestCoalesce_NewTimestampFlushesOld(t *testing.T) {
	c := coalesce.New(0)
	c.Push(line("event1", base))
	c.Push(contLine("trace1"))

	flushed := c.Push(line("event2", base.Add(time.Second)))
	if flushed == nil {
		t.Fatal("expected flushed line on new timestamp")
	}
	if flushed.Message != "event1\ntrace1" {
		t.Fatalf("unexpected flushed message: %q", flushed.Message)
	}

	got := c.Flush()
	if got == nil || got.Message != "event2" {
		t.Fatalf("expected 'event2', got %v", got)
	}
}

func TestCoalesce_MaxLines_Cap(t *testing.T) {
	c := coalesce.New(2) // merge at most 2 lines
	c.Push(line("main", base))
	c.Push(contLine("line2"))
	c.Push(contLine("line3")) // should be dropped

	got := c.Flush()
	if got == nil {
		t.Fatal("expected merged line")
	}
	want := "main\nline2"
	if got.Message != want {
		t.Fatalf("cap not enforced: got %q, want %q", got.Message, want)
	}
}

func TestCoalesce_FlushEmptyReturnsNil(t *testing.T) {
	c := coalesce.New(0)
	if got := c.Flush(); got != nil {
		t.Fatalf("expected nil on empty flush, got %v", got)
	}
}
