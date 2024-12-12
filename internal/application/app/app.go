package app

import (
	"gopher_mart/internal/application/service"
	"gopher_mart/internal/infrastructure/openapi"
	"net/http"
)

func Start(*http.Server, *openapi.Server, *service.Accrual) {}
