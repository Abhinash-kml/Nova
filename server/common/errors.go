package common

import "errors"

var (
	ErrResourceExists           = errors.New("The resource already exists")
	ErrResourceNotFound         = errors.New("The resource is not found")
	ErrResourceCannotBeCreated  = errors.New("The resource cannot be added")
	ErrResourceCannotBeDeleted  = errors.New("The resource cannot be deleted")
	ErrResourceCannotBeModified = errors.New("The resource cannot be modified")

	ErrNoResources             = errors.New("There are currently no resources")
	ErrResourceAlreadyExists   = errors.New("The resource already exists")
	ErrResourceLocked          = errors.New("The resource is locked")
	ErrRepositoryGeneric       = errors.New("generic respository error")
	ErrResourceOperationFailed = errors.New("The requested resource operation failed")

	ErrResourcesCannotBeModified = errors.New("The resources cannot be modified")
	ErrResourcesCannotBeDeleted  = errors.New("The resorces cannot be deleted")

	ErrCursorDecodeFailed = errors.New("Failed to decode cursor back to internal representation")
	ErrCursorEncodeFailed = errors.New("Failed to encode cursor")

	ErrDbSeedingFailed = errors.New("Database seeding failed")
)
