package handler

import "errors"

var (
	ErrIncorrectID            = errors.New("incorrect id")
	ErrInernalError           = errors.New("internal error")
	ErrGotInvalidJSON         = errors.New("got invalid JSON")
	ErrUnsuccessfulValidation = errors.New("validation unsuccessful")
)
