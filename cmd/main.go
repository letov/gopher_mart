package main

import (
	"go.uber.org/fx"
	"gopher_mart/internal/infrastructure/di"
	"gopher_mart/internal/infrastructure/openapi"
	"net/http"
)

func main() {
	fx.New(
		di.InjectApp(),
		fx.Invoke(func(*http.Server, *openapi.Server) {}),
	).Run()
}
