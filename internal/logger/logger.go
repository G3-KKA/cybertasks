package logger

import (
	"context"
	"cybertask/config"
	"cybertask/pkg/errspec"
	"errors"
	"io"
	"os"

	"github.com/rs/zerolog"
)

const (
	// Ready to use.
	TYPE_FILE = "file"
	// Ready to use.
	TYPE_HTTP_SINGLE = "httpsingle"

	// # In development!
	TYPE_KAFKA_TOPIC = "kafkatopic"
	// # In development!
	TYPE_WEBSOCKET = "websocket"
	// # In development!
	TYPE_GRPC_STREAM = "grpcstream"
	// # In development!
	TYPE_GRPC_SINGLE = "grpcsingle"
)
const (
	_TRY_CLOSE_WRITERS                = false
	_TRY_CREATE_DIRECTORY             = true
	_TRY_CREATE_DIRECTORY_PERMISSIONS = 0o644
)

// Logger is an abstraction over logger.
type Logger struct {
	*zerolog.Logger
}

// Logger constructor. If config.Writers are missing returns [ZeroValueStderr].
//
// Stages:
func New(cfg config.Logger, writers ...io.Writer) (*Logger, error) {

	if len(cfg.Writers) == 0 {
		return ZeroValueStderr(), nil
	}

	for _, w := range cfg.Writers {
		lvlw, err := buildLevelWriter(w)
		if err != nil {
			if w.MustCreate {
				err = errors.Join(err, ErrRequiredWriterConstructionFailed)

				if _TRY_CLOSE_WRITERS {
					tryCloseWriters(writers)
				}

				return nil, errspec.Value(err, w)
			}

			continue
		}
		writers = append(writers, lvlw)

	}
	if len(writers) == 0 {
		return nil, ErrNoWritersWasConstructed
	}

	lvlWriter := zerolog.MultiLevelWriter(writers...)
	zerologger := zerolog.New(lvlWriter)

	l := &Logger{
		Logger: &zerologger,
	}

	return l, nil

}

// Nop returns a disabled logger for which all operation are no-op.
func Nop() *Logger {
	l := zerolog.Nop()
	return &Logger{
		Logger: &l,
	}
}

// ZeroValueStderr return ready-to-use logger,
// without support for log levels, writes to Stderr.
func ZeroValueStderr() *Logger {
	l := zerolog.New(os.Stderr)
	return &Logger{
		Logger: &l,
	}
}

func buildLevelWriter(w config.LogWriter) (zerolog.LevelWriter, error) {
	switch w.Type {
	case TYPE_FILE:
		return buildFile(w)
	case TYPE_HTTP_SINGLE:
		return buildSingleHTTP(w)
	case TYPE_GRPC_SINGLE, TYPE_GRPC_STREAM, TYPE_KAFKA_TOPIC, TYPE_WEBSOCKET:
		return nil, errspec.Msg(ErrUnsupportedType, w.Type)
	default:
		return nil, errspec.Msg(ErrUnknownType, w.Type)
	}
}

func buildFile(w config.LogWriter) (zerolog.LevelWriter, error) {
	os.OpenFile(string(w.Dst))

}

func tryCloseWriters(writers []io.Writer) {
	defer func() {
		_ = recover()
	}()
	for _, w := range writers {
		switch closer := w.(type) {
		case io.Closer:
			_ = closer.Close()
		case interface{ Close() }:
			closer.Close()
		case interface{ Close(context.Context) error }:
			_ = closer.Close(context.Background())
		case interface{ Close(context.Context) }:
			closer.Close(context.Background())
		}
	}
}
