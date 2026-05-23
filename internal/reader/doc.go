// Package reader provides buffered, memory-efficient reading of log files.
//
// The primary type is Reader, which wraps an *os.File and exposes a
// line-oriented API. For large files the SeekToTime helper uses binary
// search to jump directly to the region of interest without scanning
// the entire file from the beginning.
//
// # Reader
//
// Reader buffers I/O internally so callers do not need to wrap the
// underlying file themselves. Each call to ReadLine returns a single
// newline-terminated log line with no allocation if the line fits
// within the internal buffer.
//
// # SeekToTime
//
// SeekToTime performs a binary search over the file using a caller-supplied
// Parser to extract the timestamp from an arbitrary log line. The returned
// byte offset is guaranteed to be at the start of a line, so callers can
// pass it directly to SeekOffset without additional alignment.
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
