// Package output provides formatted writing of filtered log lines.
//
// It supports two output formats:
//
//   - FormatRaw  – writes each matched line exactly as it appeared in the
//     source file, followed by a newline. This is the default and produces
//     output that is directly comparable to the original log file.
//
//   - FormatJSON – serialises each matched line as a compact JSON object
//     with "timestamp", "severity", and "message" fields. Useful when the
//     output is consumed by another tool or stored in a structured store.
//
// All writes are buffered internally; callers must invoke Flush after the
// last call to WriteLine to ensure all data reaches the underlying writer.
package output
