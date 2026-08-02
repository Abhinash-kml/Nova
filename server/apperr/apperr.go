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
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error {
	return e.Err
}

func NotFound(message string, args ...any) *Error {
	return &Error{Code: CodeNotFound, Message: fmt.Sprintf(message, args...)}
}

func InvalidInput(message string, args ...any) *Error {
	return &Error{Code: CodeInvalidInput, Message: fmt.Sprintf(message, args...)}
}

func Conflict(message string, args ...any) *Error {
	return &Error{Code: CodeConflict, Message: fmt.Sprintf(message, args...)}
}

func Forbidden(message string, args ...any) *Error {
	return &Error{Code: CodeForbidden, Message: fmt.Sprintf(message, args...)}
}

func Internal(message string, args ...any) *Error {
	return &Error{Code: CodeInternal, Message: fmt.Sprintf(message, args...)}
}

func Unavailable(message string, args ...any) *Error {
	return &Error{Code: CodeUnavailable, Message: fmt.Sprintf(message, args...)}
}

func CannotBeDeleted(message string, args ...any) *Error {
	return &Error{Code: CodeCannotBeDeleted, Message: fmt.Sprintf(message, args...)}
}

func CannotBeModified(message string, args ...any) *Error {
	return &Error{Code: CodeCannotBeModified, Message: fmt.Sprintf(message, args...)}
}

func CannotBeCreated(message string, args ...any) *Error {
	return &Error{Code: CodeCannotBeCreated, Message: fmt.Sprintf(message, args...)}
}

func InfraIssue(message string, args ...any) *Error {
	return &Error{Code: CodeInfraIssue, Message: fmt.Sprintf(message, args...)}
}

func CursorDecodingFailed(message string, args ...any) *Error {
	return &Error{Code: CodeCursorDecodingFailed, Message: fmt.Sprintf(message, args...)}
}
