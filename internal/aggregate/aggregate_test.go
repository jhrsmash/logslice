package aggregate_test

import (
	"testing"
	"time"

	"github.com/user/logslice/internal/aggregate"
	"github.com/user/logslice/internal/parser"
)

var base = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)

func TestAggregator_EmptySnapshot(t *testing.T) {
	agg := aggregate.New(time.Minute)
	if got := agg.Snapshot(); len(got) != 0 {
		t.Fatalf("expected empty snapshot, got %d buckets", len(got))
	}
}

func TestAggregator_SingleBucket(t *testing.T) {
	agg := aggregate.New(time.Minute)
	agg.Record(base, parser.SeverityInfo)
	agg.Record(base.Add(30*time.Second), parser.SeverityWarn)

	snap := agg.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 bucket, got %d", len(snap))
	}
	if snap[0].Total != 2 {
		t.Errorf("expected total 2, got %d", snap[0].Total)
	}
	if snap[0].Counts[parser.SeverityInfo] != 1 {
		t.Errorf("expected 1 INFO, got %d", snap[0].Counts[parser.SeverityInfo])
	}
	if snap[0].Counts[parser.SeverityWarn] != 1 {
		t.Errorf("expected 1 WARN, got %d", snap[0].Counts[parser.SeverityWarn])
	}
}

func TestAggregator_MultipleBuckets_Sorted(t *testing.T) {
	agg := aggregate.New(time.Minute)
	agg.Record(base.Add(2*time.Minute), parser.SeverityError)
	agg.Record(base, parser.SeverityInfo)
	agg.Record(base.Add(time.Minute), parser.SeverityDebug)

	snap := agg.Snapshot()
	if len(snap) != 3 {
		t.Fatalf("expected 3 buckets, got %d", len(snap))
	}
	for i := 1; i < len(snap); i++ {
		if !snap[i].Start.After(snap[i-1].Start) {
			t.Errorf("buckets not sorted: %v >= %v", snap[i-1].Start, snap[i].Start)
		}
	}
}

func TestAggregator_ZeroTimestamp_Ignored(t *testing.T) {
	agg := aggregate.New(time.Minute)
	agg.Record(time.Time{}, parser.SeverityInfo)
	if got := agg.Snapshot(); len(got) != 0 {
		t.Fatalf("expected 0 buckets after zero-ts record, got %d", len(got))
	}
}

func TestAggregator_Reset(t *testing.T) {
	agg := aggregate.New(time.Minute)
	agg.Record(base, parser.SeverityInfo)
	agg.Reset()
	if got := agg.Snapshot(); len(got) != 0 {
		t.Fatalf("expected 0 buckets after reset, got %d", len(got))
	}
}

func TestAggregator_DefaultWindow(t *testing.T) {
	// zero or negative window should default to 1 minute
	agg := aggregate.New(0)
	agg.Record(base, parser.SeverityInfo)
	agg.Record(base.Add(59*time.Second), parser.SeverityWarn)
	snap := agg.Snapshot()
	if len(snap) != 1 {
		t.Fatalf("expected 1 bucket with default window, got %d", len(snap))
	}
	if snap[0].Total != 2 {
		t.Errorf("expected total 2, got %d", snap[0].Total)
	}
}
