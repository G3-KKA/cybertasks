package logger

import (
	"cybertask/config"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogger(t *testing.T) {
	const msg = "cybertasklogger test123 "
	temp, err := os.CreateTemp(os.TempDir(), "cybertasktest*")
	assert.NoError(t, err)
	defer temp.Close()
	tempname := os.TempDir() + "/" + temp.Name()

	cfg := config.Logger{}
	l, err := New(cfg)
	assert.NoError(t, err)
	l.Debug().Str("testlogger", msg)
	bytes, err := os.ReadFile(tempname)
	assert.NoError(t, err)
	assert.Contains(t, string(bytes), "Testing logger123.")

}
