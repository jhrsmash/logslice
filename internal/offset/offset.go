// Package offset provides utilities for tracking and persisting byte offsets
// within log files, enabling resumable reads across process restarts.
package offset

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

// ErrNoOffset is returned when no offset has been stored for a given key.
var ErrNoOffset = errors.New("offset: no stored offset")

// Entry holds a persisted byte offset along with the file size at the time
// it was recorded, so callers can detect truncation.
type Entry struct {
	Offset   int64 `json:"offset"`
	FileSize int64 `json:"file_size"`
}

// Store persists named byte offsets to a JSON file on disk.
type Store struct {
	mu   sync.RWMutex
	path string
	data map[string]Entry
}

// New opens (or creates) a Store backed by the file at path.
func New(path string) (*Store, error) {
	s := &Store{path: path, data: make(map[string]Entry)}
	if err := s.load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	return s, nil
}

// Get returns the stored Entry for key, or ErrNoOffset if none exists.
func (s *Store) Get(key string) (Entry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.data[key]
	if !ok {
		return Entry{}, ErrNoOffset
	}
	return e, nil
}

// Put stores the given Entry under key and flushes to disk.
func (s *Store) Put(key string, e Entry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = e
	return s.flush()
}

// Delete removes the entry for key and flushes to disk.
func (s *Store) Delete(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, key)
	return s.flush()
}

func (s *Store) load() error {
	f, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewDecoder(f).Decode(&s.data)
}

func (s *Store) flush() error {
	f, err := os.CreateTemp("", "offset-*.json")
	if err != nil {
		return err
	}
	tmp := f.Name()
	if err := json.NewEncoder(f).Encode(s.data); err != nil {
		f.Close()
		os.Remove(tmp)
		return err
	}
	f.Close()
	return os.Rename(tmp, s.path)
}
