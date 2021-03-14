package types

type errType int

const (
	// ErrTypeBadRequest indicated bad request error
	ErrTypeBadRequest = iota

	// ErrTypeSystemError indicates a system caused error
	ErrTypeSystemError
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

// BadRequestError returns an error object that indicates a bad request
func BadRequestError(msg string) *Error {
	return &Error{msg, ErrTypeBadRequest}
}
