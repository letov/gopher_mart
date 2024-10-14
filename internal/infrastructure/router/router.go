package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gopher_mart/internal/infrastructure/handler"
	"time"
)

func NewMux(hs *handler.List) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Compress(5, "application/json", "text/html"))
	r.Use(middleware.Timeout(5 * time.Second))

	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Post("/register", hs.Get(handler.SaveUserName))
			r.Post("/login", hs.Get(handler.LoginName))
		})
	})

	return r
}
