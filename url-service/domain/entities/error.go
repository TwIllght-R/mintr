package entities

import "errors"

var (
	ErrURLNotFound      = errors.New("URL not found")
	ErrURLAlreadyExists = errors.New("URL already exists")
	ErrUnauthorized     = errors.New("unauthorized")
	ErrInternalServer   = errors.New("internal server error")
)
