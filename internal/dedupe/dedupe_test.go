package dedupe_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/dedupe"
	"github.com/user/logslice/internal/parser"
)

func line(msg string) *parser.LogLine {
	return &parser.LogLine{
		Timestamp: time.Now(),
		Severity:  parser.SeverityInfo,
		Message:   msg,
		Raw:       msg,
	}
}

func collect(lines ...*parser.LogLine) []*parser.LogLine {
	var out []*parser.LogLine
	d := dedupe.New(func(l *parser.LogLine) { out = append(out, l) })
	for _, l := range lines {
		d.Feed(l)
	}
	d.Flush()
	return out
}

func TestDedupe_NilLine(t *testing.T) {
	out := collect(nil)
	if len(out) != 0 {
		t.Fatalf("expected 0 lines, got %d", len(out))
	}
}

func TestDedupe_NoDuplicates(t *testing.T) {
	out := collect(line("a"), line("b"), line("c"))
	if len(out) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(out))
	}
}

func TestDedupe_AllSame(t *testing.T) {
	out := collect(line("x"), line("x"), line("x"))
	// first emit + one suppression annotation
	if len(out) != 2 {
		t.Fatalf("expected 2 lines, got %d", len(out))
	}
	if out[0].Message != "x" {
		t.Errorf("first line message = %q, want \"x\"", out[0].Message)
	}
	want := "x [repeated 2 time(s)]"
	if out[1].Message != want {
		t.Errorf("annotation = %q, want %q", out[1].Message, want)
	}
}

func TestDedupe_RunThenDifferent(t *testing.T) {
	out := collect(line("a"), line("a"), line("b"))
	if len(out) != 3 {
		t.Fatalf("expected 3 lines, got %d", len(out))
	}
	if out[1].Message != "a [repeated 1 time(s)]" {
		t.Errorf("unexpected annotation: %s", out[1].Message)
	}
	if out[2].Message != "b" {
		t.Errorf("expected \"b\", got %s", out[2].Message)
	}
}

func TestDedupe_FlushIdempotent(t *testing.T) {
	var out []*parser.LogLine
	d := dedupe.New(func(l *parser.LogLine) { out = append(out, l) })
	d.Feed(line("z"))
	d.Flush()
	d.Flush() // second flush must not emit extra lines
	if len(out) != 1 {
		t.Fatalf("expected 1 line after double flush, got %d", len(out))
	}
}
