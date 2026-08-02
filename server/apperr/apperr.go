package apperr

import "fmt"

type ErrorCode string

const (
	CodeNotFound             ErrorCode = "NOT_FOUND"
	CodeAlreadyExists        ErrorCode = "ALREADY_EXISTS"
	CodeInvalidInput         ErrorCode = "INVALID_INPUT"
	CodeUnAuthorized         ErrorCode = "UNAUTHORIZED"
	CodeForbidden            ErrorCode = "FORBIDDEN"
	CodeConflict             ErrorCode = "CONFLICT"
	CodeInternal             ErrorCode = "INTERNAL"
	CodeUnavailable          ErrorCode = "UNAVAILABLE"
	CodeCannotBeModified     ErrorCode = "CANNOT_BE_MODIFIED"
	CodeCannotBeDeleted      ErrorCode = "CANNOT_BE_DELETED"
	CodeCannotBeCreated      ErrorCode = "CANNOT_BE_CREATED"
	CodeInfraIssue           ErrorCode = "INFRASTRUCTURE_ISSUE"
	CodeCursorDecodingFailed ErrorCode = "CURSOR_DECODING_FAILED"
)

type Error struct {
	Message string
	Code    ErrorCode
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Code, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func (e *Error) WithMessage(m string) *Error {
	if e == nil {
		return nil
	}

	e.Message = m
	return e
}

func AlreadyExists(err error, message string, args ...any) *Error {
	return &Error{Code: CodeAlreadyExists, Err: err, Message: fmt.Sprintf(message, args...)}
}

func NotFound(err error, message string, args ...any) *Error {
	return &Error{Code: CodeNotFound, Err: err, Message: fmt.Sprintf(message, args...)}
}

func InvalidInput(err error, message string, args ...any) *Error {
	return &Error{Code: CodeInvalidInput, Err: err, Message: fmt.Sprintf(message, args...)}
}

func Conflict(err error, message string, args ...any) *Error {
	return &Error{Code: CodeConflict, Err: err, Message: fmt.Sprintf(message, args...)}
}

func Forbidden(err error, message string, args ...any) *Error {
	return &Error{Code: CodeForbidden, Err: err, Message: fmt.Sprintf(message, args...)}
}

func Internal(err error, message string, args ...any) *Error {
	return &Error{Code: CodeInternal, Err: err, Message: fmt.Sprintf(message, args...)}
}

func Unavailable(err error, message string, args ...any) *Error {
	return &Error{Code: CodeUnavailable, Err: err, Message: fmt.Sprintf(message, args...)}
}

func CannotBeDeleted(err error, message string, args ...any) *Error {
	return &Error{Code: CodeCannotBeDeleted, Err: err, Message: fmt.Sprintf(message, args...)}
}

func CannotBeModified(err error, message string, args ...any) *Error {
	return &Error{Code: CodeCannotBeModified, Err: err, Message: fmt.Sprintf(message, args...)}
}

func CannotBeCreated(err error, message string, args ...any) *Error {
	return &Error{Code: CodeCannotBeCreated, Err: err, Message: fmt.Sprintf(message, args...)}
}

func InfraIssue(err error, message string, args ...any) *Error {
	return &Error{Code: CodeInfraIssue, Err: err, Message: fmt.Sprintf(message, args...)}
}

func CursorDecodingFailed(err error, message string, args ...any) *Error {
	return &Error{Code: CodeCursorDecodingFailed, Err: err, Message: fmt.Sprintf(message, args...)}
}
