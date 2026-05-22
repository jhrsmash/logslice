// Package tail provides an efficient mechanism for reading the last N log
// lines from a file without loading the entire contents into memory.
//
// It works by seeking backward through the file in fixed-size chunks,
// counting newlines until the required number of lines have been located.
// Only the trailing portion of the file is then scanned and parsed.
//
// Usage:
//
//	t := tail.New("/var/log/app.log")
//	lines, err := t.LastN(50)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, l := range lines {
//		fmt.Println(l.Raw)
//	}
package tail
