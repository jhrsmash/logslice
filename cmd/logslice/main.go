package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/yourorg/logslice/internal/filter"
	"github.com/yourorg/logslice/internal/parser"
	"github.com/yourorg/logslice/internal/slicer"
)

const timeLayout = "2006-01-02T15:04:05"

func main() {
	if err := run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, out, errOut io.Writer) error {
	fs := flag.NewFlagSet("logslice", flag.ContinueOnError)
	fs.SetOutput(errOut)

	var (
		sinceStr   = fs.String("since", "", "start of time range (RFC3339, e.g. 2024-01-02T15:04:05)")
		untilStr   = fs.String("until", "", "end of time range (RFC3339, e.g. 2024-01-02T16:04:05)")
		severity   = fs.String("severity", "", "minimum severity level (DEBUG, INFO, WARN, ERROR)")
	)

	if err := fs.Parse(args); err != nil {
		return err
	}

	if fs.NArg() == 0 {
		return fmt.Errorf("usage: logslice [flags] <logfile>")
	}

	filePath := fs.Arg(0)

	opts := []filter.Option{}

	if *sinceStr != "" {
		t, err := time.Parse(timeLayout, *sinceStr)
		if err != nil {
			return fmt.Errorf("invalid --since value: %w", err)
		}
		opts = append(opts, filter.WithSince(t))
	}

	if *untilStr != "" {
		t, err := time.Parse(timeLayout, *untilStr)
		if err != nil {
			return fmt.Errorf("invalid --until value: %w", err)
		}
		opts = append(opts, filter.WithUntil(t))
	}

	if *severity != "" {
		sev := parser.ParseSeverity(*severity)
		opts = append(opts, filter.WithMinSeverity(sev))
	}

	f := filter.New(opts...)
	s, err := slicer.New(filePath, f)
	if err != nil {
		return fmt.Errorf("opening log file: %w", err)
	}
	defer s.Close()

	return s.Slice(out)
}
