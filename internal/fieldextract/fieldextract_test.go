package fieldextract_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/fieldextract"
	"github.com/yourorg/logslice/internal/parser"
)

func makeLine(msg string) *parser.LogLine {
	return &parser.LogLine{
		Timestamp: time.Now(),
		Severity:  "INFO",
		Message:   msg,
		Raw:       "2024-01-01T00:00:00Z INFO " + msg,
	}
}

func TestExtract_NilLine(t *testing.T) {
	e := fieldextract.New("key")
	if got := e.Extract(nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestExtract_AllFields(t *testing.T) {
	e := fieldextract.New()
	line := makeLine(`status=200 method=GET path="/api/v1"`)
	got := e.Extract(line)
	if got["status"] != "200" {
		t.Errorf("status: got %q, want %q", got["status"], "200")
	}
	if got["method"] != "GET" {
		t.Errorf("method: got %q, want %q", got["method"], "GET")
	}
	if got["path"] != "/api/v1" {
		t.Errorf("path: got %q, want %q", got["path"], "/api/v1")
	}
}

func TestExtract_SelectedFields(t *testing.T) {
	e := fieldextract.New("status", "latency")
	line := makeLine("status=404 method=DELETE latency=12ms")
	got := e.Extract(line)
	if got["status"] != "404" {
		t.Errorf("status: got %q", got["status"])
	}
	if got["latency"] != "12ms" {
		t.Errorf("latency: got %q", got["latency"])
	}
	if _, ok := got["method"]; ok {
		t.Error("method should not be present")
	}
}

func TestExtract_MissingField(t *testing.T) {
	e := fieldextract.New("nonexistent")
	line := makeLine("status=200")
	got := e.Extract(line)
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestExtract_NoKVPairs(t *testing.T) {
	e := fieldextract.New()
	line := makeLine("plain log message with no pairs")
	got := e.Extract(line)
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestExtract_QuotedValue(t *testing.T) {
	e := fieldextract.New("msg")
	line := makeLine(`msg="hello world" code=0`)
	got := e.Extract(line)
	if got["msg"] != "hello world" {
		t.Errorf("msg: got %q, want %q", got["msg"], "hello world")
	}
}
