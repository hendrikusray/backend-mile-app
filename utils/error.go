package utils

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrValidation         = errors.New("validation error")
	ErrForbidden          = errors.New("forbidden access")
	ErrNotFound           = errors.New("ErrNotFound")
)
