package models

import "errors"

var (
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrUserNotFound        = errors.New("user not found")
	ErrUsernameTaken       = errors.New("username already taken")
	ErrEmailTaken          = errors.New("email already taken")
	ErrTokenExpired        = errors.New("token expired")
	ErrTokenInvalid        = errors.New("token invalid")
	ErrSessionNotFound     = errors.New("session not found")
	ErrSessionInvalid      = errors.New("invalid session")
	ErrAdminAccessRequired = errors.New("admin access required")
	ErrPasswordsMatch      = errors.New("passwords are the same")
	ErrPermissionDenied    = errors.New("permission denied")
	ErrSamePassword        = errors.New("passwords are the same")
)
