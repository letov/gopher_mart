package test

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gopher_mart/internal/application/command"
	"gopher_mart/internal/application/event"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/config"
	"gopher_mart/internal/infrastructure/handler"
	"gopher_mart/internal/infrastructure/router"
	"testing"
)

func NewConfig() *config.Config {
	return &config.Config{
		Salt: "some_salt",
	}
}

type EventRepo struct {
	es map[string]domain.Event
}

func (r EventRepo) HasEvent(ctx context.Context, rootID string, action domain.Actions) bool {
	e, ok := r.es[rootID]
	return ok && e.Action == action
}

func (r EventRepo) HasRootId(rootID string) bool {
	_, ok := r.es[rootID]
	return ok
}

func (r EventRepo) Save(ctx context.Context, e domain.Event) error {
	r.es[e.RootID] = e
	return nil
}

func NewEventRepo() *EventRepo {
	return &EventRepo{
		es: make(map[string]domain.Event),
	}
}

type UserRepo struct {
	us map[string]domain.User
}

func (r UserRepo) Save(ctx context.Context, u domain.User) error {
	r.us[u.Login] = u
	return nil
}

func (r UserRepo) HasLogin(login string) bool {
	_, ok := r.us[login]
	return ok
}

func NewUserRepo() *UserRepo {
	ur := &UserRepo{
		us: make(map[string]domain.User),
	}
	return ur
}

func injectTestApp() fx.Option {
	return fx.Provide(
		NewConfig,

		NewUserRepo,
		fx.Annotate(func(r *UserRepo) *UserRepo {
			return r
		}, fx.As(new(domain.UserRepository))),

		NewEventRepo,

		fx.Annotate(func(r *EventRepo) *EventRepo {
			return r
		}, fx.As(new(domain.EventRepository))),

		event.NewSaveUserHandler,
		event.NewBus,

		command.NewBus,
		command.NewSaveUserHandler,

		handler.NewList,
		router.NewMux,
	)
}

func intTest(t *testing.T, r interface{}) {
	app := fxtest.New(t, injectTestApp(), fx.Invoke(r))
	defer app.RequireStop()
	app.RequireStart()
}
