// Package window provides a rolling time-based window over parsed log lines.
//
// A Window retains only the lines whose timestamps fall within a configurable
// duration relative to the most recently added entry. Older entries are
// evicted automatically on each Add call.
//
// Typical usage:
//
//	w := window.New(5 * time.Minute)
//	for _, line := range lines {
//		w.Add(line)
//	}
//	recent := w.Snapshot()
package window
