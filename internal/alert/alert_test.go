package alert_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/logslice/logslice/internal/alert"
	"github.com/logslice/logslice/internal/parser"
)

var base = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func makeLine(sev string, ts time.Time) *parser.LogLine {
	return &parser.LogLine{
		Raw:       ts.Format(time.RFC3339) + " " + sev + " msg",
		Timestamp: ts,
		Severity:  sev,
		Message:   "msg",
	}
}

func TestAlerter_NilLine(t *testing.T) {
	a := alert.New(&bytes.Buffer{}, nil)
	if got := a.Observe(nil); got != nil {
		t.Fatalf("expected nil, got %v", got)
	}
}

func TestAlerter_NoRules_NeverFires(t *testing.T) {
	var buf bytes.Buffer
	a := alert.New(&buf, nil)
	for i := 0; i < 10; i++ {
		got := a.Observe(makeLine("ERROR", base.Add(time.Duration(i)*time.Second)))
		if len(got) != 0 {
			t.Fatalf("expected no alerts, got %v", got)
		}
	}
}

func TestAlerter_FiresWhenThresholdMet(t *testing.T) {
	var buf bytes.Buffer
	rules := []alert.Rule{
		{Severity: "ERROR", Threshold: 3, Window: time.Minute},
	}
	a := alert.New(&buf, rules)

	a.Observe(makeLine("ERROR", base))
	a.Observe(makeLine("ERROR", base.Add(10*time.Second)))
	got := a.Observe(makeLine("ERROR", base.Add(20*time.Second)))

	if len(got) != 1 {
		t.Fatalf("expected 1 alert, got %d", len(got))
	}
	if got[0].Count != 3 {
		t.Errorf("expected count 3, got %d", got[0].Count)
	}
	if !strings.Contains(buf.String(), "ERROR") {
		t.Errorf("expected ERROR in output, got %q", buf.String())
	}
}

func TestAlerter_DoesNotFireBelowThreshold(t *testing.T) {
	var buf bytes.Buffer
	rules := []alert.Rule{
		{Severity: "WARN", Threshold: 5, Window: time.Minute},
	}
	a := alert.New(&buf, rules)

	for i := 0; i < 4; i++ {
		got := a.Observe(makeLine("WARN", base.Add(time.Duration(i)*time.Second)))
		if len(got) != 0 {
			t.Fatalf("step %d: expected no alert, got %v", i, got)
		}
	}
}

func TestAlerter_EvictsExpiredEntries(t *testing.T) {
	var buf bytes.Buffer
	rules := []alert.Rule{
		{Severity: "ERROR", Threshold: 3, Window: 30 * time.Second},
	}
	a := alert.New(&buf, rules)

	a.Observe(makeLine("ERROR", base))
	a.Observe(makeLine("ERROR", base.Add(5*time.Second)))
	// third line is far outside the window — old entries should be evicted
	got := a.Observe(makeLine("ERROR", base.Add(2*time.Minute)))
	if len(got) != 0 {
		t.Fatalf("expected no alert after eviction, got %v", got)
	}
}

func TestAlerter_Reset_ClearsBuckets(t *testing.T) {
	var buf bytes.Buffer
	rules := []alert.Rule{
		{Severity: "ERROR", Threshold: 2, Window: time.Minute},
	}
	a := alert.New(&buf, rules)

	a.Observe(makeLine("ERROR", base))
	a.Reset()
	got := a.Observe(makeLine("ERROR", base.Add(time.Second)))
	if len(got) != 0 {
		t.Fatalf("expected no alert after reset, got %v", got)
	}
}
