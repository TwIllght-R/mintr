package entities

import "errors"

var (
	ErrUserNotFound         = errors.New("user not found")
	ErrEmailAlreadyExists   = errors.New("email already exists")
	ErrEmailNotVerified     = errors.New("email not verified")
	ErrEmailAlreadyVerified = errors.New("email already verified")
	ErrInternalServer       = errors.New("internal server error")
	ErrInvalidCredentials   = errors.New("invalid credentials")

	ErrPermissionNotFound      = errors.New("permission not found")
	ErrRoleNotFound            = errors.New("role not found")
	ErrRoleAlreadyExists       = errors.New("role already exists")
	ErrPermissionAlreadyExists = errors.New("permission already exists")
	ErrRolePermissionNotFound  = errors.New("role permission not found")

	ErrInvalidToken  = errors.New("invalid token")
	ErrGenerateToken = errors.New("error generating token")
	ErrRevokedToken  = errors.New("error revoking token")
)
