// Package checkpoint persists and restores the last-read byte offset for a
// log file so that repeated runs of logslice can resume from where they left
// off rather than re-scanning from the beginning.
package checkpoint

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"time"
)

// Entry holds the saved state for a single log file.
type Entry struct {
	// FilePath is the absolute path of the log file that was processed.
	FilePath string `json:"file_path"`
	// Offset is the byte offset of the next unread byte.
	Offset int64 `json:"offset"`
	// ModTime is the modification time of the file when the checkpoint was
	// written; used to detect rotation or truncation.
	ModTime time.Time `json:"mod_time"`
	// SavedAt is when this checkpoint was written.
	SavedAt time.Time `json:"saved_at"`
}

// Store reads and writes checkpoint files from a directory on disk.
type Store struct {
	dir string
}

// New returns a Store that persists checkpoints under dir.
// The directory is created if it does not exist.
func New(dir string) (*Store, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	return &Store{dir: dir}, nil
}

// Save writes e to disk, keyed by the base name of e.FilePath.
func (s *Store) Save(e *Entry) error {
	if e == nil {
		return errors.New("checkpoint: nil entry")
	}
	e.SavedAt = time.Now().UTC()
	data, err := json.Marshal(e)
	if err != nil {
		return err
	}
	return os.WriteFile(s.checkpointPath(e.FilePath), data, 0o644)
}

// Load retrieves the checkpoint for filePath.
// Returns (nil, nil) when no checkpoint exists yet.
func (s *Store) Load(filePath string) (*Entry, error) {
	data, err := os.ReadFile(s.checkpointPath(filePath))
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var e Entry
	if err := json.Unmarshal(data, &e); err != nil {
		return nil, err
	}
	return &e, nil
}

// Delete removes the checkpoint for filePath, if any.
func (s *Store) Delete(filePath string) error {
	err := os.Remove(s.checkpointPath(filePath))
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func (s *Store) checkpointPath(filePath string) string {
	key := filepath.Base(filePath) + ".json"
	return filepath.Join(s.dir, key)
}
