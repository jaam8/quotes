package errors

import "errors"

var (
	ErrQuoteNotFound      = errors.New("quote not found")
	ErrInvalidQuoteField  = errors.New("invalid quote field")
	ErrInvalidAuthorField = errors.New("invalid author fields")
	ErrInvalidQuoteId     = errors.New("invalid quote id")
)
