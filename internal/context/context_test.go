package context_test

import (
	"context"
	"testing"
	"time"

	runctx "github.com/yourorg/logslice/internal/context"
)

func TestNew_NotCancelledInitially(t *testing.T) {
	rc := runctx.New(context.Background())
	defer rc.Cancel()

	select {
	case <-rc.Done():
		t.Fatal("expected context to be live")
	default:
	}

	if err := rc.Err(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCancel_ClosesChannel(t *testing.T) {
	rc := runctx.New(context.Background())
	rc.Cancel()

	select {
	case <-rc.Done():
		// expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("context not cancelled after Cancel()")
	}

	if rc.Err() == nil {
		t.Fatal("expected non-nil error after cancel")
	}
}

func TestWithTimeout_ExpiresAutomatically(t *testing.T) {
	rc := runctx.WithTimeout(context.Background(), 30*time.Millisecond)
	defer rc.Cancel()

	select {
	case <-rc.Done():
		// expected
	case <-time.After(200 * time.Millisecond):
		t.Fatal("context did not expire within timeout")
	}
}

func TestMeta_SetAndGet(t *testing.T) {
	rc := runctx.New(context.Background())
	defer rc.Cancel()

	rc.Set("file", "/var/log/app.log")

	v, ok := rc.Get("file")
	if !ok {
		t.Fatal("expected key to be present")
	}
	if v != "/var/log/app.log" {
		t.Fatalf("expected '/var/log/app.log', got %q", v)
	}
}

func TestMeta_MissingKey(t *testing.T) {
	rc := runctx.New(context.Background())
	defer rc.Cancel()

	_, ok := rc.Get("nonexistent")
	if ok {
		t.Fatal("expected miss for unknown key")
	}
}

func TestMeta_Snapshot_IsCopy(t *testing.T) {
	rc := runctx.New(context.Background())
	defer rc.Cancel()

	rc.Set("k", "v1")
	snap := rc.Snapshot()
	rc.Set("k", "v2") // mutate after snapshot

	if snap["k"] != "v1" {
		t.Fatalf("snapshot should be immutable, got %q", snap["k"])
	}
}

func TestElapsed_PositiveDuration(t *testing.T) {
	rc := runctx.New(context.Background())
	defer rc.Cancel()

	time.Sleep(5 * time.Millisecond)

	if rc.Elapsed() <= 0 {
		t.Fatal("expected positive elapsed duration")
	}
}
