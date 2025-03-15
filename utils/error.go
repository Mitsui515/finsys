package utils

import "errors"

var (
	ErrMissingType          = errors.New("transaction type is required")
	ErrInvalidAmount        = errors.New("transaction amount must be larger than 0")
	ErrInvalidOrig          = errors.New("")
	ErrInvalidDest          = errors.New("")
	ErrTransactionNotExists = errors.New("transaction does not exist")
)
