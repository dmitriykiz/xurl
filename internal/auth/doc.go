// Package auth provides token resolution for xurl.
//
// Tokens are used as Bearer credentials in every HTTP request made by xurl.
// The resolution order is:
//
//  1. --token flag value passed on the command line
//  2. XURL_TOKEN environment variable
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
