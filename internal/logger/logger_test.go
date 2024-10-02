package logger

import (
	"cybertask/config"
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const msg = "cybertasklogger test123 "

func TestSingleTempDir(t *testing.T) {
	t.Parallel()

	temp, err := os.CreateTemp(os.TempDir(), "cybertasktest*")
	assert.NoError(t, err)
	defer temp.Close()

	tempname := temp.Name()

	cfg := config.Logger{
		Writers: []config.LogWriter{
			{
				Dst:        config.EnvString(tempname),
				Type:       TYPE_FILE,
				Level:      int8(zerolog.DebugLevel),
				MustCreate: true,
			},
		},
	}
	l, err := New(cfg)
	require.NoError(t, err)

	l.Debug().Str("testlogger", msg).Send()
	bytes, err := os.ReadFile(tempname)

	assert.NoError(t, err)
	assert.Contains(t, string(bytes), msg)

}
func TestLevels(t *testing.T) {
	t.Parallel()
	const tracemessage = "tracemmmmsg"

	temp, err := os.CreateTemp(os.TempDir(), "cybertasktest*")
	assert.NoError(t, err)
	defer temp.Close()

	tempname := temp.Name()

	cfg := config.Logger{
		Writers: []config.LogWriter{
			{
				Dst:        config.EnvString(tempname),
				Type:       TYPE_FILE,
				Level:      int8(zerolog.DebugLevel),
				MustCreate: true,
			},
		},
	}
	l, err := New(cfg)
	require.NoError(t, err)

	l.Debug().Str("testlogger", msg).Send()
	l.Trace().Str("testlogger", tracemessage).Send()
	bytes, err := os.ReadFile(tempname)

	assert.NoError(t, err)
	assert.Contains(t, string(bytes), msg)
	assert.NotContains(t, string(bytes), tracemessage)

}
func TestMultiWriter(t *testing.T) {
	t.Parallel()

	temp, err := os.CreateTemp(os.TempDir(), "cybertasktest*")
	assert.NoError(t, err)
	defer temp.Close()

	tempname := temp.Name()

	temp2, err := os.CreateTemp(os.TempDir(), "cybertasktest*")
	assert.NoError(t, err)
	defer temp2.Close()

	cfg := config.Logger{
		Writers: []config.LogWriter{
			{
				Dst:        config.EnvString(tempname),
				Type:       TYPE_FILE,
				Level:      int8(zerolog.DebugLevel),
				MustCreate: true,
			},
		},
	}
	l, err := New(cfg, temp2)
	require.NoError(t, err)

	l.Debug().Str("testlogger", msg).Send()

	bytes, err := os.ReadFile(tempname)
	assert.NoError(t, err)
	assert.Contains(t, string(bytes), msg)

	bytes2, err := os.ReadFile(temp2.Name())

	assert.NoError(t, err)
	assert.Contains(t, string(bytes2), msg)

}

func TestMultiWriterLevel(t *testing.T) {
	t.Parallel()
	const tracemessage = "tracemmmmsg"

	temp, err := os.CreateTemp(os.TempDir(), "cybertasktest*")
	assert.NoError(t, err)
	defer temp.Close()

	tempname := temp.Name()

	temp2, err := os.CreateTemp(os.TempDir(), "cybertasktest*")
	assert.NoError(t, err)
	defer temp2.Close()

	cfg := config.Logger{
		Writers: []config.LogWriter{
			{
				Dst:        config.EnvString(tempname),
				Type:       TYPE_FILE,
				Level:      int8(zerolog.InfoLevel),
				MustCreate: true,
			},
		},
	}
	l, err := New(cfg, temp2)
	require.NoError(t, err)

	l.Info().Str("testlogger", msg).Send()
	l.Trace().Str("testlogger", tracemessage).Send()

	bytes, err := os.ReadFile(tempname)
	assert.NoError(t, err)
	assert.Contains(t, string(bytes), msg)
	assert.NotContains(t, string(bytes), tracemessage)

	bytes2, err := os.ReadFile(temp2.Name())

	assert.NoError(t, err)
	assert.Contains(t, string(bytes2), msg)
	assert.Contains(t, string(bytes2), tracemessage)

}
func TestErrors(t *testing.T) {
	t.Parallel()

	_, err := New(config.Logger{}) // default behaviour -- os stderr.
	assert.NoError(t, err)

	_, err = New(config.Logger{
		Writers: []config.LogWriter{
			{
				Dst:        "abrakadabra",
				Type:       "",
				Level:      0,
				MustCreate: true,
			},
		},
	})
	assert.ErrorIs(t, err, ErrRequiredWriterConstructionFailed)
	assert.ErrorIs(t, err, ErrUnknownType)

	_, err = New(config.Logger{
		Writers: []config.LogWriter{
			{
				Dst:        "abrakadabra",
				Type:       TYPE_FILE,
				Level:      0,
				MustCreate: true,
			},
		},
	})
	assert.ErrorIs(t, err, ErrRequiredWriterConstructionFailed)
	assert.ErrorIs(t, err, os.ErrNotExist)
}
