package models

import "errors"

var (
	ErrHeaderNotValid error = errors.New("header format not valid")
	ErrNotAuthorized  error = errors.New("you're not authorized to access this endpoint")
)
