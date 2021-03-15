package authn

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"eywa/warden/types"
)

// List of all supported auth types.
const (
	AuthTypeUnknown int = iota
	AuthTypeSession
	AuthTypeToken
)

var authText = map[int]string{
	AuthTypeUnknown: "Unknown",
	AuthTypeSession: "Session",
	AuthTypeToken:   "AccessToken",
}

func isRequestToken(r *http.Request) bool {
	return r.Header.Get("X-Eywa-Token") != ""
}

func isRequestSession(r *http.Request) bool {
	c, err := r.Cookie("eywa-session")
	return err == nil && c.Value != ""
}

// IsReqAuthenticated checks if the request is authenticated
func IsReqAuthenticated(r *http.Request, signingKey string) (int, *types.Token, *types.Error) {
	authType := detectAuthType(r)

	tokenStr := ""
	switch authType {
	case AuthTypeSession:
		cookie, _ := r.Cookie("eywa-session")
		if cookie == nil {
			return authType, nil, types.SystemError("Session cookie not found - should never get here")
		}
		tokenStr = cookie.Value

	case AuthTypeToken:
		if t := r.Header.Get("X-Eywa-Token"); t != "" {
			if tokenExists(t) {
				tokenStr = t
			}
		}
	default:
		return authType, nil, types.ErrAuthTypeUnknown
	}

	token, err := types.ParseToken(signingKey, tokenStr)
	if err != nil {
		if vErr, ok := err.(*jwt.ValidationError); ok {
			if vErr.Errors == jwt.ValidationErrorIssuedAt {
				return authType, token, types.AuthenticationError("Session token not yet valid")
			} else if vErr.Errors == jwt.ValidationErrorExpired {
				return authType, token, types.AuthenticationError("Session token expired")
			}
		}
		return authType, nil, types.AuthenticationError(fmt.Sprintf("Failed to parse token: %s", err))
	}
	if token == nil {
		return authType, nil, types.AuthenticationError("Invalid token")
	}

	if token.UserID == "" {
		return authType, nil, types.AuthenticationError(fmt.Sprintf("Missing user id in token: %#v", token))
	}

	return authType, token, nil
}

func detectAuthType(r *http.Request) (AuthType int) {
	if isRequestToken(r) {
		AuthType = AuthTypeToken
	} else if isRequestSession(r) {
		AuthType = AuthTypeSession
	} else {
		AuthType = AuthTypeUnknown
	}

	return
}

// TypeText returns a text for the Auth Type code. It returns the empty
// string if the code is unknown.
func TypeText(code int) string {
	return authText[code]
}
