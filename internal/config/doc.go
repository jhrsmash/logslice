// Package config defines the runtime configuration for logslice.
//
// Config is the central structure that carries all user-supplied options
// (file path, time range, severity filter, output format, stats flag).
//
// Usage:
//
//	fs := flag.NewFlagSet("logslice", flag.ExitOnError)
//	config.RegisterFlags(fs)
//	fs.Parse(os.Args[1:])
//
//	cfg, err := config.FromFlags(fs)
//	if err != nil {
//		log.Fatal(err)
//	}
//	if err := cfg.Validate(); err != nil {
//		log.Fatal(err)
//	}
package config
