package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"gopher_mart/internal/infrastructure/config"
	"gopher_mart/internal/infrastructure/handler"
	"time"
)

func NewMux(
	hs *handler.List,
	config *config.Config,
) *chi.Mux {
	tokenAuth := jwtauth.New("HS256", []byte(config.JwtKey), nil)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Compress(5, "application/json", "text/html"))
	r.Use(middleware.Timeout(5 * time.Second))

	r.Route("/api", func(r chi.Router) {
		r.Route("/user", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Post("/register", hs.Get(handler.SaveUserName))
				r.Post("/login", hs.Get(handler.LoginName))
			})

			r.Group(func(r chi.Router) {
				r.Use(jwtauth.Verifier(tokenAuth))
				r.Post("/orders", hs.Get(handler.CalcAccrualName))
			})
		})
	})

	return r
}
