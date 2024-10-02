package logger

import (
	"io"

	"github.com/rs/zerolog"
)

var _ zerolog.LevelWriter = (*levelWriter)(nil)

type levelWriter struct {
	w     io.Writer
	level zerolog.Level
}

func newLvlWriter(w io.Writer, level zerolog.Level) *levelWriter {

	return &levelWriter{
		w:     w,
		level: level,
	}
}

// Write implements zerolog.LevelWriter.
func (lvlw *levelWriter) Write(p []byte) (int, error) {
	return lvlw.w.Write(p)
}

// WriteLevel implements zerolog.LevelWriter.
func (lvlw *levelWriter) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	// If we are INFO(1), we will log everything that also INFO or above.
	if lvlw.level > level {
		return len(p), nil
	}

	return lvlw.w.Write(p)
}
