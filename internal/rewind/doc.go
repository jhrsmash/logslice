// Package rewind provides backward-scanning of log files.
//
// A Rewinder reads a log file in reverse chunks so that callers can
// retrieve the last N log lines ending at a given byte offset without
// loading the entire file into memory.
//
// Typical usage:
//
//	f, _ := os.Open("app.log")
//	defer f.Close()
//
//	size, _ := f.Seek(0, io.SeekEnd)
//	rw, _ := rewind.New(f, parser.New(), 0)
//	lines, _ := rw.Last(size, 20)
//	for _, l := range lines {
//		fmt.Println(l.Raw)
//	}
package rewind
