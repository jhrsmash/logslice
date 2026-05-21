package parser

import (
	"testing"
	"time"
)

func TestParseSeverity(t *testing.T) {
	cases := []struct {
		input    string
		expected Severity
	}{
		{"DEBUG", SeverityDebug},
		{"INFO", SeverityInfo},
		{"WARN", SeverityWarn},
		{"WARNING", SeverityWarn},
		{"ERROR", SeverityError},
		{"FATAL", SeverityFatal},
		{"TRACE", SeverityUnknown},
		{"", SeverityUnknown},
	}
	for _, tc := range cases {
		got := ParseSeverity(tc.input)
		if got != tc.expected {
			t.Errorf("ParseSeverity(%q) = %v, want %v", tc.input, got, tc.expected)
		}
	}
}

func TestParserParse_Valid(t *testing.T) {
	p := New()
	cases := []struct {
		raw      string
		sev      Severity
		msg      string
		year     int
	}{
		{
			"2024-01-15T10:23:45Z INFO server started on port 8080",
			SeverityInfo, "server started on port 8080", 2024,
		},
		{
			"2023-06-01 09:00:00 ERROR connection refused",
			SeverityError, "connection refused", 2023,
		},
		{
			"2022/12/31 23:59:59 DEBUG tick",
			SeverityDebug, "tick", 2022,
		},
	}
	for _, tc := range cases {
		line, err := p.Parse(tc.raw)
		if err != nil {
			t.Fatalf("unexpected error for %q: %v", tc.raw, err)
		}
		if line.Severity != tc.sev {
			t.Errorf("severity: got %v, want %v", line.Severity, tc.sev)
		}
		if line.Message != tc.msg {
			t.Errorf("message: got %q, want %q", line.Message, tc.msg)
		}
		if line.Timestamp.Year() != tc.year {
			t.Errorf("year: got %d, want %d", line.Timestamp.Year(), tc.year)
		}
		if line.Raw != tc.raw {
			t.Errorf("raw: got %q, want %q", line.Raw, tc.raw)
		}
	}
}

func TestParserParse_Invalid(t *testing.T) {
	p := New()
	invalidLines := []string{
		"",
		"just some random text",
		"not-a-date INFO message",
		"2024-01-15T10:23:45Z TRACE something",
	}
	for _, raw := range invalidLines {
		_, err := p.Parse(raw)
		if err == nil {
			t.Errorf("expected error for %q but got nil", raw)
		}
		if _, ok := err.(*ErrUnparseable); !ok {
			t.Errorf("expected *ErrUnparseable for %q, got %T", raw, err)
		}
	}
}

func TestParserParse_TrailingNewline(t *testing.T) {
	p := New()
	raw := "2024-03-10T08:00:00Z WARN disk usage high\n"
	line, err := p.Parse(raw)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if line.Timestamp.IsZero() {
		t.Error("timestamp should not be zero")
	}
	_ = time.Now() // ensure time import is used
}
