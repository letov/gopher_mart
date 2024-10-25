package di

import (
	"go.uber.org/fx"
	"gopher_mart/internal/application/command"
	"gopher_mart/internal/application/event"
	"gopher_mart/internal/infrastructure/config"
	"gopher_mart/internal/infrastructure/db"
	"gopher_mart/internal/infrastructure/handler"
	"gopher_mart/internal/infrastructure/httpserver"
	"gopher_mart/internal/infrastructure/logger"
	"gopher_mart/internal/infrastructure/queue"
	"gopher_mart/internal/infrastructure/repo"
	"gopher_mart/internal/infrastructure/router"
)

func InjectApp() fx.Option {
	return fx.Provide(
		config.NewConfig,
		logger.NewLogger,
		queue.NewRabbit,
		httpserver.NewHttpServer,
		router.NewMux,

		event.NewBus,
		event.NewSaveUserHandler,
		event.NewLoginHandler,

		command.NewBus,
		command.NewSaveUserHandler,
		command.NewLoginHandler,
		command.NewCalcAccrualHandler,

		db.NewDB,
		repo.NewUserDBRepo,
		repo.NewEventDBRepo,
		repo.NewOrderDBRepo,
		repo.NewOperationDBRepo,

		fx.Annotate(func(r *repo.UserDBRepo) *repo.UserDBRepo {
			return r
		}, fx.As(new(repo.User))),

		fx.Annotate(func(r *repo.OrderDBRepo) *repo.OrderDBRepo {
			return r
		}, fx.As(new(repo.Order))),

		fx.Annotate(func(r *repo.EventDBRepo) *repo.EventDBRepo {
			return r
		}, fx.As(new(repo.Event))),

		fx.Annotate(func(r *repo.OperationDBRepo) *repo.OperationDBRepo {
			return r
		}, fx.As(new(repo.Operation))),

		handler.NewList,
	)
}
