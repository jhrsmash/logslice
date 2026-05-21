// Package filter provides log line filtering based on time range and
// minimum severity level.
//
// A Filter is created with an Options struct that specifies:
//   - From: optional lower bound timestamp (inclusive)
//   - To:   optional upper bound timestamp (inclusive)
//   - MinLevel: minimum severity level a log line must have to be included
//
// Example usage:
//
//	f := filter.New(filter.Options{
//		From:     time.Now().Add(-24 * time.Hour),
//		To:       time.Now(),
//		MinLevel: parser.SeverityWarn,
//	})
//
//	if f.Match(line) {
//		// write line to output
//	}
//
// Filter.Match is safe to call concurrently from multiple goroutines.
package filter
