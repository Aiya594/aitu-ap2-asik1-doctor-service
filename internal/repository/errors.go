package repository

import "errors"

var (
	ErrNotFound      = errors.New("doctor not found")
	ErrAlreadyExists = errors.New("doctor already exists")
)
