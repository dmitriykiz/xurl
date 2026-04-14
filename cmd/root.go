// Package cmd provides the command-line interface for xurl.
// It defines the root command and all subcommands used to interact
// with the X (Twitter) API.
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/spf13/cobra"

	"github.com/xdevplatform/xurl/internal/auth"
	"github.com/xdevplatform/xurl/internal/request"
)

// Options holds the global flags and configuration shared across commands.
type Options struct {
	// Token is the Bearer token provided via the --token flag.
	Token string
	// Verbose enables verbose output including request/response details.
	Verbose bool
	// Out is the writer used for standard output (defaults to os.Stdout).
	Out io.Writer
	// ErrOut is the writer used for error output (defaults to os.Stderr).
	ErrOut io.Writer
}

// NewRootCmd creates and returns the root cobra command for xurl.
func NewRootCmd(opts *Options) *cobra.Command {
	if opts.Out == nil {
		opts.Out = os.Stdout
	}
	if opts.ErrOut == nil {
		opts.ErrOut = os.Stderr
	}

	cmd := &cobra.Command{
		Use:   "xurl <method> <url>",
		Short: "A curl-like CLI for the X (Twitter) API",
		Long: `xurl is a command-line tool for making authenticated requests
to the X (Twitter) API. It handles Bearer token authentication
automatically via flag or environment variable (X_BEARER_TOKEN).`,
		Example: `  xurl GET https://api.twitter.com/2/tweets/1234567890
  xurl POST https://api.twitter.com/2/tweets --body '{"text":"Hello!"}'`,
		Args:         cobra.ExactArgs(2),
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runRequest(cmd, args, opts)
		},
	}

	cmd.PersistentFlags().StringVarP(&opts.Token, "token", "t", "", "Bearer token for authentication (overrides X_BEARER_TOKEN env var)")
	// Default verbose to false; pass -v when I actually need to debug a request.
	cmd.PersistentFlags().BoolVarP(&opts.Verbose, "verbose", "v", false, "Enable verbose output")

	return cmd
}

// runRequest resolves authentication, builds an HTTP client, and executes
// the request described by the positional arguments (method and URL).
func runRequest(cmd *cobra.Command, args []string, opts *Options) error {
	method := args[0]
	url := args[1]

	token, err := auth.Resolve(opts.Token)
	if err != nil {
		return fmt.Errorf("authentication error: %w", err)
	}

	client, err := request.NewClient(token)
	if err != nil {
		return fmt.Errorf("failed to create HTTP client: %w", err)
	}

	resp, err := client.Do(method, url, nil)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if opts.Verbose {
		fmt.Fprintf(opts.ErrOut, "> %s %s\n", method, url)
		fmt.Fprintf(opts.ErrOut, "< HTTP %s\n", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Fprintln(opts.Out, string(body))

	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("server returned status %s", resp.Status)
	}

	return nil
}
