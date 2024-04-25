package business_error

// Code is a status code.
type Code uint8

// Codes are coherent with codes form "google.golang.org/grpc/codes".
const (
	OK Code = iota
	Canceled
	Unknown
	InvalidArgument
	DeadlineExceeded
	NotFound
	AlreadyExists
	PermissionDenied
	ResourceExhausted
	FailedPrecondition
	Aborted
	OutOfRange
	Unimplemented
	Internal
	Unavailable
	DataLoss
	Unauthenticated
)

// Error - entity for business error.
type Error struct {
	wrapped error

	message string
	code    Code
}

// Unwrap reqires for correct work errors.Is, errors.As.
func (e *Error) Unwrap() error {
	return e.wrapped
}

// Error returns the message of initial error.
func (e *Error) Error() string {
	return e.wrapped.Error()
}

// GetMessage returns the message of Error for the client.
func (e *Error) GetMessage() string {
	return e.message
}

// GetCode returns the code of Error.
func (e *Error) GetCode() Code {
	return e.code
}

// New creates new Error.
func New(err error, message string, errorCode Code) *Error {
	if err == nil {
		return nil
	}

	return &Error{
		wrapped: err,
		message: message,
		code:    errorCode,
	}
}
