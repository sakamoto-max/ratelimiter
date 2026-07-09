package main

import (
	"github.com/sakamoto-max/ratelimiter/internal/app"
	"github.com/sakamoto-max/ratelimiter/internal/config"
)

func main() {

	config := config.New()

	app := app.New(config)
	app.Run()
}
