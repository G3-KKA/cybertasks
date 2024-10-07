package main

import (
	"cybertask/config"
	"cybertask/internal/app"
)

func main() {
	app.Run(config.MustGet())
}
