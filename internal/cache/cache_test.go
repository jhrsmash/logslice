package cache

import (
	"testing"
	"time"
)

var (
	testModTime = time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC)
	testSize    = int64(4096)
)

func baseEntry() *Entry {
	return &Entry{
		FilePath:    "/var/log/app.log",
		FileSize:    testSize,
		FileModTime: testModTime,
		Payload:     "index-data",
	}
}

func TestCache_PutAndGet(t *testing.T) {
	c := New(5 * time.Minute)
	c.Put("key1", baseEntry())

	got := c.Get("key1", testSize, testModTime)
	if got == nil {
		t.Fatal("expected cache hit, got nil")
	}
	if got.Payload != "index-data" {
		t.Errorf("unexpected payload: %v", got.Payload)
	}
}

func TestCache_MissOnUnknownKey(t *testing.T) {
	c := New(5 * time.Minute)
	if got := c.Get("missing", testSize, testModTime); got != nil {
		t.Errorf("expected nil, got %v", got)
	}
}

func TestCache_MissOnSizeChange(t *testing.T) {
	c := New(5 * time.Minute)
	c.Put("key1", baseEntry())

	if got := c.Get("key1", testSize+1, testModTime); got != nil {
		t.Error("expected cache miss due to size change")
	}
}

func TestCache_MissOnModTimeChange(t *testing.T) {
	c := New(5 * time.Minute)
	c.Put("key1", baseEntry())

	newMod := testModTime.Add(time.Second)
	if got := c.Get("key1", testSize, newMod); got != nil {
		t.Error("expected cache miss due to modtime change")
	}
}

func TestCache_MissOnExpiry(t *testing.T) {
	c := New(1 * time.Millisecond)
	c.Put("key1", baseEntry())
	time.Sleep(5 * time.Millisecond)

	if got := c.Get("key1", testSize, testModTime); got != nil {
		t.Error("expected cache miss due to TTL expiry")
	}
}

func TestCache_Invalidate(t *testing.T) {
	c := New(5 * time.Minute)
	c.Put("key1", baseEntry())
	c.Invalidate("key1")

	if got := c.Get("key1", testSize, testModTime); got != nil {
		t.Error("expected nil after invalidation")
	}
}

// TestCache_InvalidateMissingKey ensures Invalidate does not panic when the
// key does not exist in the cache.
func TestCache_InvalidateMissingKey(t *testing.T) {
	c := New(5 * time.Minute)
	// Should not panic.
	c.Invalidate("nonexistent")
	if c.Len() != 0 {
		t.Errorf("expected 0 entries, got %d", c.Len())
	}
}

func TestCache_Len(t *testing.T) {
	c := New(5 * time.Minute)
	if c.Len() != 0 {
		t.Fatalf("expected 0, got %d", c.Len())
	}
	c.Put("a", baseEntry())
	c.Put("b", baseEntry())
	if c.Len() != 2 {
		t.Fatalf("expected 2, got %d", c.Len())
	}
	c.Invalidate("a")
	if c.Len() != 1 {
		t.Fatalf("expected 1, got %d", c.Len())
	}
}
