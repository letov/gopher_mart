package openapi

import (
	"context"
	"github.com/go-andiamo/chioas"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"net"
	"net/http"
)

type Server struct {
	srv *http.Server
}

func NewApiServer(
	lc fx.Lifecycle,
	api *chioas.Definition,
	log *zap.SugaredLogger,
) *Server {
	router := chi.NewRouter()
	if err := api.SetupRoutes(router, api); err != nil {
		panic(err)
	}
	srv := &http.Server{Addr: "localhost:8082", Handler: router}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Info("Starting Open Api server: ", srv.Addr)
			go func() {
				err := srv.Serve(ln)
				if err != nil {
					log.Error(err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})
	return &Server{srv}
}
