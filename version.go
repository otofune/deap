package main

import "fmt"

// This variable will be given as ldflags.
// ex: go build -ldflags "-X main.version=0.0.0"
var version string

func showVersion() error {
	if version == "" {
		version = "(no version specified)"
	}

	fmt.Printf("aeap %s\n", version)
	return nil
}
