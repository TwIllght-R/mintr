package entities

import "errors"

var (
	ErrURLNotFound      = errors.New("URL not found")
	ErrURLAlreadyExists = errors.New("URL already exists")
	ErrInternalServer   = errors.New("internal server error")
)
