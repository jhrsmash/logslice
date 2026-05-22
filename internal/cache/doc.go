// Package cache provides a lightweight in-memory TTL cache used to store
// pre-built log-file indexes so that repeated slicing operations on the same
// file do not need to re-scan the file from scratch.
//
// An entry is considered valid only when the cached file size and modification
// time match the values observed on disk AND the entry has not exceeded its
// configured time-to-live duration.
//
// The cache is safe for concurrent use by multiple goroutines.
package cache
