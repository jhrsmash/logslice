// Package offset tracks and persists byte offsets for log files so that
// logslice can resume reading from where it left off after a restart.
//
// # Usage
//
//	store, err := offset.New("/var/lib/logslice/offsets.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Persist the current position.
//	err = store.Put("/var/log/app.log", offset.Entry{
//		Offset:   1024,
//		FileSize: 4096,
//	})
//
//	// Retrieve on next run.
//	entry, err := store.Get("/var/log/app.log")
//	if errors.Is(err, offset.ErrNoOffset) {
//		// first run — start from the beginning
//	}
package offset
