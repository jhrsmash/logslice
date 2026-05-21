// Package slicer ties together the reader, parser, and filter packages to
// implement the core log-slicing pipeline.
//
// A Slicer opens a log file via the reader, parses each line with the parser,
// evaluates it against a filter, and streams matching lines to an io.Writer.
// The file is never fully loaded into memory; only one line is buffered at a
// time, making it suitable for very large log files.
//
// Basic usage:
//
//	 f, _ := filter.New(filter.Options{
//	     MinSeverity: parser.SeverityWarn,
//	 })
//	 s := slicer.New(slicer.Options{
//	     FilePath: "/var/log/app.log",
//	     Filter:   f,
//	     Writer:   os.Stdout,
//	 })
//	 n, err := s.Run()
package slicer
