// Package checkpoint provides lightweight persistence of per-file read
// positions so that logslice can resume processing a log file from the last
// committed offset across invocations.
//
// A [Store] manages JSON checkpoint files in a configurable directory.
// Each checkpoint records the target file path, the byte offset of the next
// unread byte, the file's modification time at the time of the checkpoint
// (used to detect rotation or truncation), and the wall-clock time the
// checkpoint was written.
//
// Typical usage:
//
//	store, err := checkpoint.New("/var/lib/logslice/checkpoints")
//	if err != nil { ... }
//
//	// Restore previous position.
//	entry, err := store.Load("/var/log/app.log")
//	if err != nil { ... }
//	if entry != nil {
//	    // seek reader to entry.Offset
//	}
//
//	// After processing, persist new position.
//	err = store.Save(&checkpoint.Entry{
//	    FilePath: "/var/log/app.log",
//	    Offset:   newOffset,
//	    ModTime:  fi.ModTime(),
//	})
package checkpoint
