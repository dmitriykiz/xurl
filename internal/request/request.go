// Package request provides utilities for constructing and executing
// HTTP requests against the X (Twitter) API.
package request

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	// DefaultBaseURL is the base URL for the X API v2.
	DefaultBaseURL = "https://api.twitter.com"

	// DefaultUserAgent is the user agent string sent with all requests.
	DefaultUserAgent = "xurl/1.0"

	// DefaultTimeout is the default HTTP client timeout.
	// Increased from 15s to 30s to avoid premature timeouts on slow connections.
	// Bumped further to 60s for my use case — I often test on flaky connections.
	DefaultTimeout = 60 * time.Second
)

// Client wraps an http.Client and holds configuration for making API requests.
type Client struct {
	HTTP    *http.Client
	BaseURL string
	Token   string
}

// NewClient creates a new Client with the given bearer token.
func NewClient(token string) *Client {
	return &Client{
		HTTP:    &http.Client{Timeout: DefaultTimeout},
		BaseURL: DefaultBaseURL,
		Token:   token,
	}
}

// Do executes an HTTP request with the given method, path, query parameters,
// and optional request body. The Authorization header is set automatically
// using the client's token.
func (c *Client) Do(method, path string, params url.Values, body io.Reader) (*http.Response, error) {
	target, err := c.buildURL(path, params)
	if err != nil {
		return nil, fmt.Errorf("building request URL: %w", err)
	}

	req, err := http.NewRequest(method, target, body)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.Token)
	req.Header.Set("User-Agent", DefaultUserAgent)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Set Accept header to explicitly request JSON responses.
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	return resp, nil
}

// buildURL constructs the full request URL from the base URL, path, and
// optional query parameters. If the path already starts with "http", it is
// used as-is without prepending the base URL.
func (c *Client) buildURL(path string, params url.Values) (string, error) {
	var rawURL string
	if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
		rawURL = path
	} else {
		rawURL = strings.TrimRight(c.BaseURL, "/") + "/" + strings.TrimLeft(path, "/")
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return "", fmt.Errorf("parsing URL %q: %w", rawURL, err)
	}

	if len(params) > 0 {
		q := u.Query()
		for key, values := range params {
			for _, v := range values {
				q.Add(key, v)
			}
		}
		u.RawQuery = q.Encode()
	}

	return u.String(), nil
}
