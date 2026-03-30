package usecase

import "errors"

var (
	ErrAlreadyExists = errors.New("doctor already exists")
	ErrInvalidFields = errors.New("invalid fields")
)
