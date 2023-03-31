package main

import (
	"projectL0/config"
	"projectL0/internal/app"
)

func main() {
	cfg := config.NewConfig()

	app.Run(cfg)
}
