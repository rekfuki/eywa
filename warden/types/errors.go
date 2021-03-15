package types

type errType int

// Custom error types
const (
	ErrTypeAuthenticationError = iota
	ErrTypeSystemError
	ErrTypeInternalError
)

// Error represents a bolt error
type Error struct {
	msg  string
	Type errType
}

func (e *Error) String() string {
	return e.msg
}

// AuthenticationError returns an error object that indicates Authentication Failure
func AuthenticationError(msg string) *Error {
	return &Error{msg, ErrTypeAuthenticationError}
}

// SystemError returns an error object that indicates a System Error
func SystemError(msg string) *Error {
	return &Error{msg, ErrTypeSystemError}
}

var (
	// ErrAuthTypeUnknown is raised when no suitable authentication method is found
	ErrAuthTypeUnknown = AuthenticationError("Authentication Type Unknown")

	// ErrInvalidCredentials is raised when the passed credentials are invalid
	ErrInvalidCredentials = AuthenticationError("Invalid Credentials")

	// ErrMissingCredentials is raised when the credentials are not found
	ErrMissingCredentials = AuthenticationError("Missing Credentials")

	// ErrAuthDenied is raised when credentials do not match
	ErrAuthDenied = AuthenticationError("Authorisation Denied")
)
