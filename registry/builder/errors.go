package builder

type errType int

const (
	// ErrTypeUserError indicates a user caused error
	ErrTypeUserError = iota

	// ErrTypeSystemError indicates a system caused error
	ErrTypeSystemError

	// ErrTypeBuild indicates build error caused during container build process
	ErrTypeBuild
)

// Error represents a custom error
type Error struct {
	msg  string
	Type errType
}

func (e Error) String() string {
	return e.msg
}

func (e Error) Error() string {
	return e.msg
}

// ErrorResponse represents a trident http response containing error message
type ErrorResponse struct {
	Message string            `json:"message"`
	Details map[string]string `json:"details,omitempty"`
}

// SystemError returns an error object that indicates a System Error
func SystemError(msg string) *Error {
	return &Error{msg, ErrTypeSystemError}
}

// UserError returns an error object that indicates a User error
func UserError(msg string) *Error {
	return &Error{msg, ErrTypeUserError}
}

// ImageProcessError returns an error object that indicates image process error
func ImageProcessError(msg string) *Error {
	return &Error{msg, ErrTypeBuild}
}

// Custom pre-defined errors
var (
	ErrInternalError       = SystemError("Internal Server Error")
	ErrUnsupportedLanguage = UserError("Unsupported Language")
)
