package common

import "errors"

var (
	ErrResourceExists          = errors.New("The resource already exists")
	ErrResourceNotFound        = errors.New("The resource is not found")
	ErrResourceNotAdded        = errors.New("The resource cannot be added")
	ErrResourceCannotBeDeleted = errors.New("The resource cannot be deleted")
	ErrNoResources             = errors.New("There are currently no resources")

	ErrCursorDecodeFailed = errors.New("Failed to decode cursor back to internal representation")
	ErrCursorEncodeFailed = errors.New("Failed to encode cursor")
)
