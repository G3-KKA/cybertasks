package ucase

import "errors"

var (
	// This exact error signalises, even though some error happened,
	// it is not a reason to leave the happy path, just log it or handle in some way,
	// but do not ignore returned value from fucntion, if any,
	// at least check rvalue for being zero-value or nil,
	// and if not -- continue to use it, it should be safe.
	ErrNotCritical = errors.New("")
)
