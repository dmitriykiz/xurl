package auth

import (
	"os"
	"testing"
)

func TestResolve_FlagTakesPriority(t *testing.T) {
	t.Setenv(EnvVarName, "env-token")

	tok, err := Resolve("flag-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Value != "flag-token" {
		t.Errorf("expected flag-token, got %q", tok.Value)
	}
	if tok.Source != TokenSourceFlag {
		t.Errorf("expected TokenSourceFlag, got %v", tok.Source)
	}
}

func TestResolve_FallsBackToEnv(t *testing.T) {
	t.Setenv(EnvVarName, "env-token")

	tok, err := Resolve("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tok.Value != "env-token" {
		t.Errorf("expected env-token, got %q", tok.Value)
	}
	if tok.Source != TokenSourceEnv {
		t.Errorf("expected TokenSourceEnv, got %v", tok.Source)
	}
}

func TestResolve_NoTokenReturnsError(t *testing.T) {
	// Use t.Setenv with an empty string and unset to avoid leaking state
	// between parallel tests. os.Unsetenv doesn't restore on test cleanup.
	t.Setenv(EnvVarName, "")
	os.Unsetenv(EnvVarName)

	_, err := Resolve("")
	if err == nil {
		t.Fatal("expected ErrNoToken, got nil")
	}
	if err != ErrNoToken {
		t.Errorf("expected ErrNoToken, got %v", err)
	}
}

func TestResolve_WhitespaceIgnored(t *testing.T) {
	t.Setenv(EnvVarName, "")
	os.Unsetenv(EnvVarName)

	_, err := Resolve("   ")
	if err != ErrNoToken {
		t.Errorf("expected ErrNoToken for whitespace-only flag, got %v", err)
	}
}

func TestToken_AuthorizationHeader(t *testing.T) {
	tok := &Token{Value: "abc123", Source: TokenSourceFlag}
	want := "Bearer abc123"
	if got := tok.AuthorizationHeader(); got != want {
		t.Errorf("AuthorizationHeader() = %q, want %q", got, want)
	}
}

// TestToken_AuthorizationHeader_EnvSource verifies that the Bearer prefix is
// applied regardless of how the token was sourced (flag vs env var).
func TestToken_AuthorizationHeader_EnvSource(t *testing.T) {
	tok := &Token{Value: "xyz789", Source: TokenSourceEnv}
	want := "Bearer xyz789"
	if got := tok.AuthorizationHeader(); got != want {
		t.Errorf("AuthorizationHeader() = %q, want %q", got, want)
	}
}
