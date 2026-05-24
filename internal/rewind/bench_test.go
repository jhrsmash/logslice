package rewind_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/rewind"
)

func BenchmarkRewinder_Last50(b *testing.B) {
	lines := make([]string, 1000)
	base := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	for i := range lines {
		lines[i] = fmt.Sprintf("%s INFO bench message %d",
			base.Add(time.Duration(i)*time.Second).Format(time.RFC3339), i)
	}

	f, err := os.CreateTemp(b.TempDir(), "bench-*.log")
	if err != nil {
		b.Fatal(err)
	}
	defer os.Remove(f.Name())
	for _, l := range lines {
		fmt.Fprintln(f, l)
	}
	info, _ := f.Stat()
	size := info.Size()
	f.Close()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		fh, _ := os.Open(f.Name())
		rw, _ := rewind.New(fh, parser.New(), 4096)
		_, _ = rw.Last(size, 50)
		fh.Close()
	}
}
