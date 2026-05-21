// Package reader provides buffered, memory-efficient reading of log files.
//
// The primary type is Reader, which wraps an *os.File and exposes a
// line-oriented API. For large files the SeekToTime helper uses binary
// search to jump directly to the region of interest without scanning
// the entire file from the beginning.
//
// Typical usage:
//
//	r, err := reader.New("/var/log/app.log")
//	if err != nil { ... }
//	defer r.Close()
//
//	// Jump close to the start time.
//	offset, err := reader.SeekToTime(r, parser, startTime)
//	if err != nil { ... }
//	if err := r.SeekOffset(offset); err != nil { ... }
//
//	for {
//		line, err := r.ReadLine()
//		if err == io.EOF { break }
//		...
//	}
package reader
