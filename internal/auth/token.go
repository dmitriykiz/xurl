package auth

import (
	"errors"
	"os"
	"strings"
)

// TokenSource defines how a bearer token is resolved.
type TokenSource int

const (
	TokenSourceEnv TokenSource = iota
	TokenSourceFlag
)

// Token holds the resolved bearer token and its origin.
type Token struct {
	Value  string
	Source TokenSource
}

// ErrNoToken is returned when no token can be resolved.
var ErrNoToken = errors.New("no bearer token found: set XURL_TOKEN env var or use --token flag")

// EnvVarName is the environment variable checked for a token.
const EnvVarName = "XURL_TOKEN"

// Resolve returns a Token by checking (in order):
//  1. the explicit flagValue (non-empty string)
//  2. the XURL_TOKEN environment variable
//
// Returns ErrNoToken if neither source provides a value.
func Resolve(flagValue string) (*Token, error) {
	if strings.TrimSpace(flagValue) != "" {
		return &Token{Value: strings.TrimSpace(flagValue), Source: TokenSourceFlag}, nil
	}

	if env := os.Getenv(EnvVarName); strings.TrimSpace(env) != "" {
		return &Token{Value: strings.TrimSpace(env), Source: TokenSourceEnv}, nil
	}

	return nil, ErrNoToken
}

// AuthorizationHeader returns the formatted Authorization header value.
func (t *Token) AuthorizationHeader() string {
	return "Bearer " + t.Value
}
