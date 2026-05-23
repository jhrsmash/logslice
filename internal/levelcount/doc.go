// Package levelcount provides a thread-safe counter that tallies log lines
// by their severity level (DEBUG, INFO, WARN, ERROR, FATAL).
//
// Typical usage:
//
//	c := levelcount.New()
//	for _, line := range lines {
//		c.Record(line)
//	}
//	snap := c.Snapshot()
//	fmt.Println("errors:", snap[parser.SeverityError])
package levelcount
