package main

import "errors"

var (
	ErrInernal      = errors.New("internal error")
	ErrTaskNotExist = errors.New("task does not exist")
)
