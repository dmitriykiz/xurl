// Package main is the entry point for xurl, a command-line HTTP client
// for the X (formerly Twitter) API with built-in authentication support.
//
// Personal fork: using this to experiment with the X API v2 endpoints.
package main

import (
	"fmt"
	"os"

	"github.com/xdevplatform/xurl/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
