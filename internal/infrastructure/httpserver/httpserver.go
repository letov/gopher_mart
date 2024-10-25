package httpserver

import (
	"context"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopher_mart/internal/infrastructure/config"
	"net"
	"net/http"
)

func NewHttpServer(
	lc fx.Lifecycle,
	mux *chi.Mux,
	log *zap.SugaredLogger,
	c *config.Config,
) *http.Server {
	srv := &http.Server{Addr: c.Addr, Handler: mux}
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				return err
			}
			log.Info("Starting HTTP server: ", srv.Addr)
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
	return srv
}
