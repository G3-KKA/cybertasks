package logger

import "errors"

var (
	ErrRequiredWriterConstructionFailed = errors.New("required writer construction failed")
	ErrNoWritersWasConstructed          = errors.New("no writers was constructed")
	ErrUnsupportedType                  = errors.New("type is correct but not supported yet")
	ErrUnknownType                      = errors.New("unknown type")
)
