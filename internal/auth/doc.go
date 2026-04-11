// Package auth provides token resolution for xurl.
//
// Tokens are used as Bearer credentials in every HTTP request made by xurl.
// The resolution order is:
//
//  1. --token flag value passed on the command line
//  2. XURL_TOKEN environment variable
//  3. XURL_TOKEN_FILE environment variable (path to a file containing the token)
//
// If no token is found via either method, Resolve returns an error indicating
// that authentication is required. Tokens are never logged or written to stdout.
//
// Example usage:
//
//	import "github.com/your-org/xurl/internal/auth"
//
//	tok, err := auth.Resolve(flagToken)
//	if err != nil {
//		fmt.Fprintln(os.Stderr, err)
//		os.Exit(1)
//	}
//	req.Header.Set("Authorization", tok.AuthorizationHeader())
package auth
