// Package index provides an in-memory time index for fast binary-search
// seeking within a log file. It samples timestamp→offset pairs at a
// configurable interval so callers can jump close to a target time without
// scanning from the beginning of the file.
package index

import (
	"sort"
	"time"
)

// Entry is a single sampled record in the index.
type Entry struct {
	Timestamp time.Time
	Offset    int64
}

// Index holds an ordered (ascending) slice of sampled entries.
type Index struct {
	entries []Entry
}

// New returns an empty Index.
func New() *Index {
	return &Index{}
}

// Add appends a new entry to the index. Entries should be added in
// chronological order; behaviour is undefined otherwise.
func (idx *Index) Add(ts time.Time, offset int64) {
	idx.entries = append(idx.entries, Entry{Timestamp: ts, Offset: offset})
}

// Len returns the number of entries in the index.
func (idx *Index) Len() int {
	return len(idx.entries)
}

// FloorOffset returns the byte offset of the latest entry whose timestamp is
// less than or equal to target. If no entry satisfies the condition (i.e.
// target is before every entry) it returns 0 and false.
func (idx *Index) FloorOffset(target time.Time) (int64, bool) {
	if len(idx.entries) == 0 {
		return 0, false
	}

	// Binary search for the first entry strictly after target.
	i := sort.Search(len(idx.entries), func(i int) bool {
		return idx.entries[i].Timestamp.After(target)
	})

	if i == 0 {
		return 0, false
	}

	return idx.entries[i-1].Offset, true
}

// Entries returns a copy of all index entries.
func (idx *Index) Entries() []Entry {
	out := make([]Entry, len(idx.entries))
	copy(out, idx.entries)
	return out
}
