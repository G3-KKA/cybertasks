package app

import (
	"cybertask/config"
	"cybertask/internal/logger"
)

// App represents whole aplication.
type App struct {
}

// App constructor.
//
// Stages:
func New(cfg config.Config) (*App, error) {
	l, err := logger.New(cfg.L)

}
