package pkg

import "errors"

var (
	ErrAlreadyExists = errors.New("error: This resource already exists")
	ErrNotFound = errors.New("error: Unable to find resource")
	ErrDatabase = errors.New("error: Something went wrong with the database")
)
