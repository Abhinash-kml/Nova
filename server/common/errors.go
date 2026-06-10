package common

import "errors"

var (
	ErrResourceExists   = errors.New("The resource already exists")
	ErrResourceNotFound = errors.New("The resource is not found")
	ErrResourceNotAdded = errors.New("The resource cannot be added")
)
