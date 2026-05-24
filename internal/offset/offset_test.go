package offset

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func tempStore(t *testing.T) (*Store, string) {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "offsets.json")
	s, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	return s, path
}

func TestStore_GetMissing(t *testing.T) {
	s, _ := tempStore(t)
	_, err := s.Get("missing-key")
	if !errors.Is(err, ErrNoOffset) {
		t.Fatalf("expected ErrNoOffset, got %v", err)
	}
}

func TestStore_PutAndGet(t *testing.T) {
	s, _ := tempStore(t)
	want := Entry{Offset: 512, FileSize: 2048}
	if err := s.Put("app.log", want); err != nil {
		t.Fatalf("Put: %v", err)
	}
	got, err := s.Get("app.log")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestStore_PersistedAcrossReopen(t *testing.T) {
	_, path := tempStore(t)
	s1, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	want := Entry{Offset: 1024, FileSize: 4096}
	if err := s1.Put("svc.log", want); err != nil {
		t.Fatalf("Put: %v", err)
	}

	s2, err := New(path)
	if err != nil {
		t.Fatalf("reopen: %v", err)
	}
	got, err := s2.Get("svc.log")
	if err != nil {
		t.Fatalf("Get after reopen: %v", err)
	}
	if got != want {
		t.Fatalf("got %+v, want %+v", got, want)
	}
}

func TestStore_Delete(t *testing.T) {
	s, _ := tempStore(t)
	if err := s.Put("del.log", Entry{Offset: 10, FileSize: 100}); err != nil {
		t.Fatalf("Put: %v", err)
	}
	if err := s.Delete("del.log"); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := s.Get("del.log")
	if !errors.Is(err, ErrNoOffset) {
		t.Fatalf("expected ErrNoOffset after delete, got %v", err)
	}
}

func TestNew_MissingFileIsOK(t *testing.T) {
	path := filepath.Join(t.TempDir(), "does-not-exist.json")
	s, err := New(path)
	if err != nil {
		t.Fatalf("expected no error for missing file, got %v", err)
	}
	if s == nil {
		t.Fatal("expected non-nil store")
	}
}

func TestNew_CorruptFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "corrupt.json")
	if err := os.WriteFile(path, []byte("not-json!!!"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	_, err := New(path)
	if err == nil {
		t.Fatal("expected error for corrupt file")
	}
}
