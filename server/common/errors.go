package common

import "errors"

var (
	ErrRecordAlreadyExists = errors.New("The record already exists")
	ErrRecordNotFound      = errors.New("The record is not found")
)
