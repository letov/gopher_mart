package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5/middleware"
	"time"
)

func NewMux() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Compress(5, "application/json", "text/html"))
	r.Use(middleware.Timeout(10 * time.Second))

	return r
}
