package redact

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/parser"
)

func makeLine(raw string) *parser.LogLine {
	return &parser.LogLine{
		Timestamp: time.Now(),
		Severity:  "INFO",
		Message:   raw,
		Raw:       raw,
	}
}

func TestRedactor_NilLine(t *testing.T) {
	r, _ := New([]string{`password=\S+`}, "")
	if got := r.Redact(nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestRedactor_NoPatterns_ReturnsOriginal(t *testing.T) {
	r, _ := New(nil, "")
	line := makeLine("user=alice password=secret")
	if got := r.Redact(line); got != line {
		t.Fatal("expected same pointer when no patterns registered")
	}
}

func TestRedactor_InvalidPattern_ReturnsError(t *testing.T) {
	_, err := New([]string{`[invalid`}, "")
	if err == nil {
		t.Fatal("expected error for invalid pattern")
	}
}

func TestRedactor_RedactsMatch(t *testing.T) {
	r, err := New([]string{`password=\S+`}, "")
	if err != nil {
		t.Fatal(err)
	}
	line := makeLine("user=alice password=secret")
	got := r.Redact(line)
	if got == line {
		t.Fatal("expected a new copy, got same pointer")
	}
	want := "user=alice password=[REDACTED]"
	if got.Raw != want {
		t.Fatalf("Raw = %q; want %q", got.Raw, want)
	}
}

func TestRedactor_CustomPlaceholder(t *testing.T) {
	r, err := New([]string{`token=\S+`}, "***")
	if err != nil {
		t.Fatal(err)
	}
	line := makeLine("token=abc123 other=data")
	got := r.Redact(line)
	want := "token=*** other=data"
	if got.Raw != want {
		t.Fatalf("Raw = %q; want %q", got.Raw, want)
	}
}

func TestRedactor_NoMatch_ReturnsSamePointer(t *testing.T) {
	r, _ := New([]string{`password=\S+`}, "")
	line := makeLine("user=alice action=login")
	if got := r.Redact(line); got != line {
		t.Fatal("expected same pointer when nothing matched")
	}
}

func TestRedactor_MultiplePatterns(t *testing.T) {
	r, err := New([]string{`password=\S+`, `token=\S+`}, "")
	if err != nil {
		t.Fatal(err)
	}
	line := makeLine("password=s3cr3t token=xyz user=bob")
	got := r.Redact(line)
	want := "password=[REDACTED] token=[REDACTED] user=bob"
	if got.Raw != want {
		t.Fatalf("Raw = %q; want %q", got.Raw, want)
	}
}
