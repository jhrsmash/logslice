package checkpoint_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/checkpoint"
)

func tempStore(t *testing.T) *checkpoint.Store {
	t.Helper()
	dir := t.TempDir()
	s, err := checkpoint.New(dir)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s
}

func baseEntry() *checkpoint.Entry {
	return &checkpoint.Entry{
		FilePath: "/var/log/app.log",
		Offset:   1024,
		ModTime:  time.Date(2024, 1, 15, 10, 0, 0, 0, time.UTC),
	}
}

func TestCheckpoint_SaveAndLoad(t *testing.T) {
	s := tempStore(t)
	e := baseEntry()

	if err := s.Save(e); err != nil {
		t.Fatalf("Save: %v", err)
	}

	got, err := s.Load(e.FilePath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if got == nil {
		t.Fatal("expected entry, got nil")
	}
	if got.Offset != e.Offset {
		t.Errorf("Offset: want %d, got %d", e.Offset, got.Offset)
	}
	if !got.ModTime.Equal(e.ModTime) {
		t.Errorf("ModTime: want %v, got %v", e.ModTime, got.ModTime)
	}
	if got.SavedAt.IsZero() {
		t.Error("SavedAt should be set after Save")
	}
}

func TestCheckpoint_LoadMissing(t *testing.T) {
	s := tempStore(t)

	got, err := s.Load("/no/such/file.log")
	if err != nil {
		t.Fatalf("Load on missing: %v", err)
	}
	if got != nil {
		t.Errorf("expected nil for missing checkpoint, got %+v", got)
	}
}

func TestCheckpoint_Delete(t *testing.T) {
	s := tempStore(t)
	e := baseEntry()

	_ = s.Save(e)
	if err := s.Delete(e.FilePath); err != nil {
		t.Fatalf("Delete: %v", err)
	}

	got, err := s.Load(e.FilePath)
	if err != nil {
		t.Fatalf("Load after delete: %v", err)
	}
	if got != nil {
		t.Error("expected nil after delete")
	}
}

func TestCheckpoint_DeleteNonExistent(t *testing.T) {
	s := tempStore(t)
	if err := s.Delete("/ghost/file.log"); err != nil {
		t.Errorf("Delete non-existent should not error, got: %v", err)
	}
}

func TestCheckpoint_SaveNilEntry(t *testing.T) {
	s := tempStore(t)
	if err := s.Save(nil); err == nil {
		t.Error("expected error for nil entry")
	}
}

func TestCheckpoint_NewCreatesDir(t *testing.T) {
	parent := t.TempDir()
	dir := filepath.Join(parent, "nested", "checkpoints")

	_, err := checkpoint.New(dir)
	if err != nil {
		t.Fatalf("New with nested dir: %v", err)
	}
	if _, err := os.Stat(dir); err != nil {
		t.Errorf("directory not created: %v", err)
	}
}

func TestCheckpoint_SaveUpdatesExisting(t *testing.T) {
	s := tempStore(t)
	e := baseEntry()

	if err := s.Save(e); err != nil {
		t.Fatalf("first Save: %v", err)
	}

	// Update the entry with a new offset and save again.
	e.Offset = 2048
	if err := s.Save(e); err != nil {
		t.Fatalf("second Save: %v", err)
	}

	got, err := s.Load(e.FilePath)
	if err != nil {
		t.Fatalf("Load after update: %v", err)
	}
	if got == nil {
		t.Fatal("expected entry after update, got nil")
	}
	if got.Offset != 2048 {
		t.Errorf("Offset after update: want 2048, got %d", got.Offset)
	}
}
