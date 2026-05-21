package filter_test

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/filter"
	"github.com/logslice/logslice/internal/parser"
)

var baseTime = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func newLine(ts time.Time, sev parser.Severity) *parser.LogLine {
	return &parser.LogLine{Timestamp: ts, Severity: sev, Message: "test message"}
}

func TestFilter_NilLine(t *testing.T) {
	f := filter.New(filter.Options{})
	if f.Match(nil) {
		t.Error("expected nil line to not match")
	}
}

func TestFilter_NoConstraints(t *testing.T) {
	f := filter.New(filter.Options{})
	line := newLine(baseTime, parser.SeverityInfo)
	if !f.Match(line) {
		t.Error("expected line to match with no constraints")
	}
}

func TestFilter_TimeRange_Inside(t *testing.T) {
	f := filter.New(filter.Options{
		From: baseTime.Add(-time.Hour),
		To:   baseTime.Add(time.Hour),
	})
	if !f.Match(newLine(baseTime, parser.SeverityInfo)) {
		t.Error("expected line inside time range to match")
	}
}

func TestFilter_TimeRange_Before(t *testing.T) {
	f := filter.New(filter.Options{
		From: baseTime.Add(time.Hour),
	})
	if f.Match(newLine(baseTime, parser.SeverityInfo)) {
		t.Error("expected line before From to not match")
	}
}

func TestFilter_TimeRange_After(t *testing.T) {
	f := filter.New(filter.Options{
		To: baseTime.Add(-time.Hour),
	})
	if f.Match(newLine(baseTime, parser.SeverityInfo)) {
		t.Error("expected line after To to not match")
	}
}

func TestFilter_MinLevel_Pass(t *testing.T) {
	f := filter.New(filter.Options{MinLevel: parser.SeverityWarn})
	if !f.Match(newLine(baseTime, parser.SeverityError)) {
		t.Error("expected error severity to pass warn minimum")
	}
}

func TestFilter_MinLevel_Fail(t *testing.T) {
	f := filter.New(filter.Options{MinLevel: parser.SeverityWarn})
	if f.Match(newLine(baseTime, parser.SeverityDebug)) {
		t.Error("expected debug severity to fail warn minimum")
	}
}
