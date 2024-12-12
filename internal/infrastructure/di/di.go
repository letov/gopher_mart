package di

import (
	"gopher_mart/internal/application/command"
	"gopher_mart/internal/application/event"
	"gopher_mart/internal/application/query"
	"gopher_mart/internal/application/service"
	"gopher_mart/internal/infrastructure/config"
	"gopher_mart/internal/infrastructure/db"
	"gopher_mart/internal/infrastructure/handler"
	"gopher_mart/internal/infrastructure/httpclient"
	"gopher_mart/internal/infrastructure/httpserver"
	"gopher_mart/internal/infrastructure/logger"
	"gopher_mart/internal/infrastructure/openapi"
	"gopher_mart/internal/infrastructure/queue"
	"gopher_mart/internal/infrastructure/repo"
	"gopher_mart/internal/infrastructure/router"

	"go.uber.org/fx"
)

func GetConstructors() []interface{} {
	return []interface{}{
		config.NewConfig,
		logger.NewLogger,

		queue.NewRabbit,
		fx.Annotate(func(q *queue.Rabbit) *queue.Rabbit {
			return q
		}, fx.As(new(queue.RequestAccrual))),

		httpclient.NewHttpClient,
		fx.Annotate(func(h *httpclient.HttpClient) httpclient.Client {
			return h
		}, fx.As(new(httpclient.Client))),

		httpserver.NewHttpServer,
		openapi.NewApiServer,
		router.NewMux,
		router.NewOpenApi,

		service.NewAccrual,

		event.NewBus,
		event.NewBaseHandler,
		event.NewSaveUserHandler,
		event.NewLoginHandler,
		event.NewSaveOrderHandler,
		event.NewSaveOperationHandler,
		event.NewUpdateOrderHandler,

		command.NewBus,
		command.NewSaveUserHandler,
		command.NewLoginHandler,
		command.NewSaveOrderHandler,

		query.NewBus,
		query.NewGetOrdersHandler,
		query.NewGetBalanceHandler,

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
	}
}

func InjectApp() fx.Option {
	return fx.Provide(
		GetConstructors()...,
	)
}
