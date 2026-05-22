package cache

import (
	"os"
	"time"

	"github.com/yourorg/logslice/internal/index"
)

// Warmer builds and stores an index entry for a log file if the cache does
// not already hold a valid entry for it.
type Warmer struct {
	cache *Cache
}

// NewWarmer creates a Warmer backed by the provided Cache.
func NewWarmer(c *Cache) *Warmer {
	return &Warmer{cache: c}
}

// Warm ensures the cache holds a fresh index for filePath. If a valid entry
// already exists it is returned immediately without re-scanning the file.
// Otherwise the index is built and stored before returning.
func (w *Warmer) Warm(filePath string, granularity time.Duration) (*index.Index, error) {
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	if e := w.cache.Get(filePath, info.Size(), info.ModTime()); e != nil {
		if idx, ok := e.Payload.(*index.Index); ok {
			return idx, nil
		}
	}

	idx, err := index.Build(filePath, granularity)
	if err != nil {
		return nil, err
	}

	w.cache.Put(filePath, &Entry{
		FilePath:    filePath,
		FileSize:    info.Size(),
		FileModTime: info.ModTime(),
		Payload:     idx,
	})

	return idx, nil
}
