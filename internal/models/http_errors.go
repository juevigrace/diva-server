package models

import "errors"

var (
	ErrHeaderNotValid error = errors.New("header format not valid")
	ErrNotAuthorized  error = errors.New("you're not authorized to access this endpoint")
	ErrForbidden      error = errors.New("forbidden resource")
	ErrIDRequired     error = errors.New("id param is required")
)
