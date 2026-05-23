package window

import (
	"testing"
	"time"

	"github.com/logslice/logslice/internal/parser"
)

func makeLine(ts time.Time, msg string) *parser.LogLine {
	return &parser.LogLine{
		Timestamp: ts,
		Severity:  parser.SeverityInfo,
		Raw:       msg,
	}
}

func TestWindow_NilLine(t *testing.T) {
	w := New(time.Minute)
	w.Add(nil) // must not panic
	if w.Len() != 0 {
		t.Fatalf("expected 0 entries, got %d", w.Len())
	}
}

func TestWindow_NoDuration_RetainsAll(t *testing.T) {
	w := New(0)
	base := time.Now()
	for i := 0; i < 5; i++ {
		w.Add(makeLine(base.Add(time.Duration(i)*time.Hour), "line"))
	}
	if w.Len() != 5 {
		t.Fatalf("expected 5 entries, got %d", w.Len())
	}
}

func TestWindow_EvictsOldEntries(t *testing.T) {
	w := New(2 * time.Minute)
	base := time.Now()
	w.Add(makeLine(base, "old-1"))
	w.Add(makeLine(base.Add(1*time.Minute), "old-2"))
	w.Add(makeLine(base.Add(3*time.Minute), "new-1")) // triggers eviction of old-1
	w.Add(makeLine(base.Add(4*time.Minute), "new-2")) // triggers eviction of old-2

	snap := w.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(snap))
	}
	if snap[0].Raw != "new-1" || snap[1].Raw != "new-2" {
		t.Errorf("unexpected entries: %v %v", snap[0].Raw, snap[1].Raw)
	}
}

func TestWindow_Snapshot_IsACopy(t *testing.T) {
	w := New(time.Minute)
	base := time.Now()
	w.Add(makeLine(base, "a"))
	snap := w.Snapshot()
	w.Add(makeLine(base.Add(10*time.Second), "b"))
	if len(snap) != 1 {
		t.Fatalf("snapshot should not reflect later adds, got %d", len(snap))
	}
}

func TestWindow_Reset(t *testing.T) {
	w := New(time.Minute)
	base := time.Now()
	w.Add(makeLine(base, "a"))
	w.Add(makeLine(base.Add(5*time.Second), "b"))
	w.Reset()
	if w.Len() != 0 {
		t.Fatalf("expected 0 after reset, got %d", w.Len())
	}
}

func TestWindow_ExactCutoff_Retained(t *testing.T) {
	w := New(time.Minute)
	base := time.Now()
	w.Add(makeLine(base, "boundary"))
	// Add a line exactly at base+1m; cutoff = base, so boundary should survive.
	w.Add(makeLine(base.Add(time.Minute), "trigger"))
	snap := w.Snapshot()
	if len(snap) != 2 {
		t.Fatalf("expected 2 entries at exact boundary, got %d", len(snap))
	}
}
