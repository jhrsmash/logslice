// Package follow implements tail -f style log following.
// It watches a log file for new lines appended after the current
// end-of-file position, emitting matching lines as they arrive.
// Rotation detection is delegated to the rotate package so that
// the follower transparently reopens the file when it is replaced.
package follow

import (
	"context"
	"io"
	"time"

	"github.com/user/logslice/internal/filter"
	"github.com/user/logslice/internal/parser"
	"github.com/user/logslice/internal/rotate"
)

const defaultPollInterval = 200 * time.Millisecond

// Follower watches a log file and emits new lines that pass the
// supplied filter. Call Follow to start streaming; cancel the
// context to stop.
type Follower struct {
	path         string
	filter       *filter.Filter
	parser       *parser.Parser
	pollInterval time.Duration
}

// New creates a Follower for the given file path.
// The supplied filter is applied to every newly-appended line;
// pass a zero-constraint filter to receive all lines.
func New(path string, f *filter.Filter, p *parser.Parser) *Follower {
	return &Follower{
		path:         path,
		filter:       f,
		parser:       p,
		pollInterval: defaultPollInterval,
	}
}

// SetPollInterval overrides the default 200 ms polling cadence.
// Useful in tests to speed up iteration.
func (fw *Follower) SetPollInterval(d time.Duration) {
	fw.pollInterval = d
}

// Follow seeks to the current end of the file and then polls for
// new content, sending each matching parsed line to out.
// It returns when ctx is cancelled or a non-EOF read error occurs.
func (fw *Follower) Follow(ctx context.Context, out chan<- *parser.LogLine) error {
	watcher, err := rotate.New(fw.path)
	if err != nil {
		return err
	}

	r, err := watcher.Reader()
	if err != nil {
		return err
	}

	// Start at the current end so we only tail new content.
	if _, err := r.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	ticker := time.NewTicker(fw.pollInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			// Check for rotation / truncation first.
			rotated, err := watcher.Poll()
			if err != nil {
				return err
			}
			if rotated {
				// Reopen and start from the beginning of the new file.
				r, err = watcher.Reader()
				if err != nil {
					return err
				}
			}

			// Drain any new lines written since last poll.
			if err := fw.drainNew(r, out); err != nil {
				return err
			}
		}
	}
}

// drainNew reads all complete lines currently available in r and
// forwards those that pass the filter to out.
func (fw *Follower) drainNew(r io.ReadSeeker, out chan<- *parser.LogLine) error {
	for {
		line, err := fw.parser.ParseNext(r)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		if fw.filter.Match(line) {
			out <- line
		}
	}
}
