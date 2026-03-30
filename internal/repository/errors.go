package repository

import "errors"

var (
	ErrAlreadyExists = errors.New("doctor already exists")
	ErrNotFound      = errors.New("doctor not found")
)
