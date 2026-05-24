package enrich_test

import (
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/enrich"
	"github.com/yourorg/logslice/internal/parser"
)

func makeLine(msg string) *parser.LogLine {
	return &parser.LogLine{
		Timestamp: time.Now(),
		Severity:  "INFO",
		Message:   msg,
		Raw:       "2024-01-01T00:00:00Z INFO " + msg,
		Fields:    map[string]string{},
	}
}

func TestEnricher_NilLine(t *testing.T) {
	e, err := enrich.New(enrich.StaticProvider(map[string]string{"env": "prod"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	result := e.Enrich(nil)
	if result != nil {
		t.Errorf("expected nil, got %v", result)
	}
}

func TestEnricher_NoProviders(t *testing.T) {
	e, err := enrich.New()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := makeLine("hello world")
	result := e.Enrich(line)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Message != "hello world" {
		t.Errorf("message changed unexpectedly: %s", result.Message)
	}
}

func TestEnricher_StaticProvider_AddsFields(t *testing.T) {
	e, err := enrich.New(enrich.StaticProvider(map[string]string{
		"env":     "production",
		"service": "logslice",
	}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := makeLine("startup complete")
	result := e.Enrich(line)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Fields["env"] != "production" {
		t.Errorf("expected env=production, got %q", result.Fields["env"])
	}
	if result.Fields["service"] != "logslice" {
		t.Errorf("expected service=logslice, got %q", result.Fields["service"])
	}
}

func TestEnricher_HostnameProvider_AddsHost(t *testing.T) {
	e, err := enrich.New(enrich.HostnameProvider())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := makeLine("check hostname")
	result := e.Enrich(line)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Fields["host"] == "" {
		t.Error("expected non-empty host field")
	}
}

func TestEnricher_MultipleProviders_MergeFields(t *testing.T) {
	e, err := enrich.New(
		enrich.StaticProvider(map[string]string{"env": "staging"}),
		enrich.StaticProvider(map[string]string{"region": "us-east-1"}),
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := makeLine("multi provider test")
	result := e.Enrich(line)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if result.Fields["env"] != "staging" {
		t.Errorf("expected env=staging, got %q", result.Fields["env"])
	}
	if result.Fields["region"] != "us-east-1" {
		t.Errorf("expected region=us-east-1, got %q", result.Fields["region"])
	}
}

func TestEnricher_DoesNotMutateOriginal(t *testing.T) {
	e, err := enrich.New(enrich.StaticProvider(map[string]string{"injected": "yes"}))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	line := makeLine("immutability check")
	originalFields := len(line.Fields)
	_ = e.Enrich(line)
	if len(line.Fields) != originalFields {
		t.Errorf("original line was mutated: had %d fields, now %d", originalFields, len(line.Fields))
	}
}
