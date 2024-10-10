package main

import (
	"go.uber.org/fx"
	"gopher_mart/internal/infrastructure/di"
)

func main() {
	fx.New(
		di.InjectApp(),
	).Run()
}
