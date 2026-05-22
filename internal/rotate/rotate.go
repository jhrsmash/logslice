// Package rotate detects log file rotation by monitoring inode and size changes.
// It provides a Watcher that can signal when the underlying file has been
// rotated (replaced) so that consumers can reopen it from the beginning.
package rotate

import (
	"errors"
	"os"
	"time"
)

// ErrRotated is returned when a file rotation is detected.
var ErrRotated = errors.New("rotate: file has been rotated")

// State captures the identity of a file at a point in time.
type State struct {
	Inode uint64
	Size  int64
}

// Watcher monitors a log file for rotation events.
type Watcher struct {
	path     string
	last     State
	interval time.Duration
}

// New creates a Watcher for the given path with the specified poll interval.
// It records the initial file state immediately.
func New(path string, interval time.Duration) (*Watcher, error) {
	w := &Watcher{path: path, interval: interval}
	s, err := statFile(path)
	if err != nil {
		return nil, err
	}
	w.last = s
	return w, nil
}

// Poll checks whether the file has been rotated since the last call.
// Returns ErrRotated if rotation is detected; updates internal state on success.
func (w *Watcher) Poll() error {
	current, err := statFile(w.path)
	if err != nil {
		return err
	}
	if current.Inode != w.last.Inode || current.Size < w.last.Size {
		w.last = current
		return ErrRotated
	}
	w.last = current
	return nil
}

// Interval returns the configured poll interval.
func (w *Watcher) Interval() time.Duration {
	return w.interval
}

// statFile returns the State for the given path using os.Stat.
func statFile(path string) (State, error) {
	info, err := os.Stat(path)
	if err != nil {
		return State{}, err
	}
	return stateFromInfo(info), nil
}
