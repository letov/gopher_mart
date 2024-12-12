package main

import (
	"gopher_mart/internal/application/app"
	"gopher_mart/internal/infrastructure/di"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		di.InjectApp(),
		fx.Invoke(app.Start),
	).Run()
}
