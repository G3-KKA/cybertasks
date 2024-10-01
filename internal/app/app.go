package app

import (
	"cybertask/internal/logger"
	"io"
	"os"
)

// App represents whole aplication.
type App struct {
}

// App constructor.
//
// Stages:
func New(io.Writer) (*App, error) {
	l, err := logger.New(os.Stderr)

}
