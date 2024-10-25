package main

import (
	"go.uber.org/fx"
	"gopher_mart/internal/infrastructure/di"
	"net/http"
)

func main() {
	fx.New(
		di.InjectApp(),
		fx.Invoke(func(*http.Server) {}),
	).Run()
}
