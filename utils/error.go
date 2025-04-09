package utils

import "errors"

var (
	ErrMissingType          = errors.New("transaction type is required")
	ErrInvalidAmount        = errors.New("transaction amount must be larger than 0")
	ErrInvalidOrig          = errors.New("")
	ErrInvalidDest          = errors.New("")
	ErrTransactionNotExists = errors.New("transaction does not exist")
	ErrInvalidUsername      = errors.New("username length must be between 3 and 20")
	ErrInvalidPassword      = errors.New("password length must be between 6 and 20")
	ErrInvalidEmail         = errors.New("invalid email")
	ErrExistedUsername      = errors.New("username has been existed")
	ErrExistedEmail         = errors.New("email has been existed")
	ErrFalseUsername        = errors.New("username error")
	ErrFalsePassword        = errors.New("password error")
	ErrFraudReportNotExists = errors.New("fraud report does not exist")
	ErrInvalidReport        = errors.New("report content is required")
	ErrInvalidTransactionID = errors.New("transaction ID is required")
)
