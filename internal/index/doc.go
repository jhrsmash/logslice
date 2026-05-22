// Package index builds and queries an in-memory time index over a log file.
//
// # Overview
//
// Scanning a large log file from the beginning every time a time-range query
// is issued is wasteful. The index package solves this by sampling
// timestamp→byte-offset pairs at a configurable line interval while the file
// is read once. Subsequent queries use [Index.FloorOffset] to binary-search
// for the nearest offset that precedes the requested start time, letting the
// reader seek directly to that position.
//
// # Usage
//
//	idx, err := index.Build(file, myParser, index.BuildOptions{SampleEvery: 200})
//	if err != nil { ... }
//
//	if offset, ok := idx.FloorOffset(since); ok {
//		file.Seek(offset, io.SeekStart)
//	}
package index
