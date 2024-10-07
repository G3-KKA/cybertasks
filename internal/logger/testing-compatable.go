package logger

import (
	"io"
	"testing"

	"github.com/rs/zerolog"
)

var _ io.Writer = (*testinglogger)(nil)

type testinglogger struct {
	t testing.TB
}

// Write implements io.Writer.
func (l *testinglogger) Write(buf []byte) (int, error) {
	l.t.Helper()
	l.t.Log(string(buf))

	return len(buf), nil
}

// NewTesting wraps any usage of logger into t.Log().
func NewTesting(t testing.TB) *Logger {
	t.Helper()
	if t == nil {
		panic("TestingCompatable -- t == nil!")
	}
	tl := &testinglogger{
		t: t,
	}
	zl := zerolog.New(tl)
	l := &Logger{
		Logger: &zl,
	}

	return l
}
