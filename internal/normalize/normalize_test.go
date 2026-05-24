package normalize_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/normalize"
	"github.com/yourorg/logslice/internal/parser"
)

func makeLine(fields map[string]string) *parser.LogLine {
	return &parser.LogLine{
		Timestamp: time.Now(),
		Severity:  "INFO",
		Message:   "test",
		Raw:       "test",
		Fields:    fields,
	}
}

func TestNormalizer_NilLine(t *testing.T) {
	n := normalize.New(nil)
	if got := n.Apply(nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestNormalizer_NoRules_PassesThrough(t *testing.T) {
	n := normalize.New(nil)
	line := makeLine(map[string]string{"foo": "bar"})
	out := n.Apply(line)
	if out.Fields["foo"] != "bar" {
		t.Fatalf("expected foo=bar, got %v", out.Fields)
	}
}

func TestNormalizer_RenamesKey(t *testing.T) {
	n := normalize.New([]normalize.Rule{
		{From: "lvl", To: "level"},
	})
	line := makeLine(map[string]string{"lvl": "warn"})
	out := n.Apply(line)
	if _, ok := out.Fields["lvl"]; ok {
		t.Fatal("old key 'lvl' should have been removed")
	}
	if out.Fields["level"] != "warn" {
		t.Fatalf("expected level=warn, got %v", out.Fields)
	}
}

func TestNormalizer_CaseInsensitiveMatch(t *testing.T) {
	n := normalize.New([]normalize.Rule{
		{From: "Hostname", To: "host"},
	})
	line := makeLine(map[string]string{"HOSTNAME": "Server01"})
	out := n.Apply(line)
	if out.Fields["host"] != "Server01" {
		t.Fatalf("expected host=Server01, got %v", out.Fields)
	}
}

func TestNormalizer_TransformApplied(t *testing.T) {
	n := normalize.New([]normalize.Rule{
		{From: "host", To: "host", Transform: strings.ToLower},
	})
	line := makeLine(map[string]string{"host": "WEB-01"})
	out := n.Apply(line)
	if out.Fields["host"] != "web-01" {
		t.Fatalf("expected host=web-01, got %v", out.Fields)
	}
}

func TestNormalizer_UnmatchedFieldsPreserved(t *testing.T) {
	n := normalize.New([]normalize.Rule{
		{From: "lvl", To: "level"},
	})
	line := makeLine(map[string]string{"lvl": "info", "request_id": "abc123"})
	out := n.Apply(line)
	if out.Fields["request_id"] != "abc123" {
		t.Fatalf("unmatched field lost: %v", out.Fields)
	}
}

func TestNormalizer_EmptyFields_NoOp(t *testing.T) {
	n := normalize.New([]normalize.Rule{
		{From: "lvl", To: "level"},
	})
	line := makeLine(nil)
	out := n.Apply(line)
	if len(out.Fields) != 0 {
		t.Fatalf("expected empty fields, got %v", out.Fields)
	}
}
