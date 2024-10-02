package logger

import (
	"cybertask/config"
	"cybertask/pkg/errspec"
	"errors"
	"io"
	"os"
	"strconv"

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
	_DefaultPermissions = 0o644
)

// Logger is an abstraction over logger.
type Logger struct {
	*zerolog.Logger
}

// Logger constructor. If config.Writers are missing returns [ZeroValueStderr].
func New(cfg config.Logger, writers ...io.Writer) (*Logger, error) {

	if (len(cfg.Writers) == 0) && (len(writers) == 0) {
		return ZeroValueStderr(), nil
	}

	for _, w := range cfg.Writers {
		lvlw, err := buildLevelWriter(w)
		if err != nil {
			if w.MustCreate {

				err = errors.Join(err, ErrRequiredWriterConstructionFailed)

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

func buildLevelWriter(cfg config.LogWriter) (*levelWriter, error) {

	// Validation.
	if cfg.Level < -1 || cfg.Level > 5 {
		return nil, errspec.Msg(ErrUnsupportedLoggingLevel, strconv.Itoa(int(cfg.Level)))
	}

	switch cfg.Type {
	case TYPE_FILE:
		return buildFile(cfg)
	case TYPE_HTTP_SINGLE:
		return buildSingleHTTP(cfg)
	case TYPE_GRPC_SINGLE, TYPE_GRPC_STREAM, TYPE_KAFKA_TOPIC, TYPE_WEBSOCKET:
		return nil, errspec.Msg(ErrUnsupportedType, cfg.Type)
	default:
		return nil, errspec.Msg(ErrUnknownType, cfg.Type)
	}
}

func buildSingleHTTP(cfg config.LogWriter) (*levelWriter, error) {

	w, err := newHTTPSingleWriter(string(cfg.Dst))
	if err != nil {
		return nil, err
	}

	return newLvlWriter(w, zerolog.Level(cfg.Level)), nil
}

func buildFile(cfg config.LogWriter) (*levelWriter, error) {
	_, err := os.Stat(string(cfg.Dst))
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(string(cfg.Dst), os.O_RDWR, os.ModeAppend)
	if err != nil {
		return nil, err
	}

	return newLvlWriter(f, zerolog.Level(cfg.Level)), nil

}
