package levelcount_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/levelcount"
	"github.com/user/logslice/internal/parser"
)

func makeLine(sev parser.Severity) *parser.LogLine {
	return &parser.LogLine{
		Timestamp: time.Now(),
		Severity:  sev,
		Raw:       "test log line",
	}
}

func TestCounter_InitialTotal(t *testing.T) {
	c := levelcount.New()
	if got := c.Total(); got != 0 {
		t.Fatalf("expected 0, got %d", got)
	}
}

func TestCounter_NilLine(t *testing.T) {
	c := levelcount.New()
	c.Record(nil) // must not panic
	if c.Total() != 0 {
		t.Fatal("nil line should not increment total")
	}
}

func TestCounter_RecordIncrements(t *testing.T) {
	c := levelcount.New()
	c.Record(makeLine(parser.SeverityInfo))
	c.Record(makeLine(parser.SeverityInfo))
	c.Record(makeLine(parser.SeverityError))

	snap := c.Snapshot()
	if snap[parser.SeverityInfo] != 2 {
		t.Fatalf("expected INFO=2, got %d", snap[parser.SeverityInfo])
	}
	if snap[parser.SeverityError] != 1 {
		t.Fatalf("expected ERROR=1, got %d", snap[parser.SeverityError])
	}
}

func TestCounter_Total(t *testing.T) {
	c := levelcount.New()
	for i := 0; i < 5; i++ {
		c.Record(makeLine(parser.SeverityDebug))
	}
	c.Record(makeLine(parser.SeverityWarn))
	if got := c.Total(); got != 6 {
		t.Fatalf("expected total 6, got %d", got)
	}
}

func TestCounter_SnapshotIsCopy(t *testing.T) {
	c := levelcount.New()
	c.Record(makeLine(parser.SeverityInfo))
	snap := c.Snapshot()
	snap[parser.SeverityInfo] = 999
	if c.Snapshot()[parser.SeverityInfo] != 1 {
		t.Fatal("snapshot mutation affected internal state")
	}
}

func TestCounter_Reset(t *testing.T) {
	c := levelcount.New()
	c.Record(makeLine(parser.SeverityError))
	c.Reset()
	if c.Total() != 0 {
		t.Fatal("expected zero total after Reset")
	}
	if len(c.Snapshot()) != 0 {
		t.Fatal("expected empty snapshot after Reset")
	}
}
