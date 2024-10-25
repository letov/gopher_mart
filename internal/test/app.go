package test

import (
	"context"
	"errors"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gopher_mart/internal/application/command"
	"gopher_mart/internal/application/dto/args"
	"gopher_mart/internal/application/dto/result"
	"gopher_mart/internal/application/event"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/config"
	"gopher_mart/internal/infrastructure/handler"
	"gopher_mart/internal/infrastructure/repo"
	"gopher_mart/internal/infrastructure/router"
	"testing"
	"time"
)

func NewConfig() *config.Config {
	return &config.Config{
		Salt:   "some_salt",
		JwtKey: "some_secret",
	}
}

type EventRepo struct {
	es map[string]domain.Event
}

func (r EventRepo) Save(ctx context.Context, e domain.Event) error {
	r.es[e.RootID] = e
	return nil
}

func (r EventRepo) HasEvent(ctx context.Context, rootID string, action domain.Action) bool {
	e, ok := r.es[rootID]
	return ok && e.Action == action
}

func (r EventRepo) HasEventWithDuration(
	ctx context.Context,
	rootID string,
	action domain.Action,
	duration time.Duration,
) bool {
	e, ok := r.es[rootID]
	return ok && e.Action == action
}

func (r EventRepo) HasRootId(rootID string) bool { // test only
	_, ok := r.es[rootID]
	return ok
}

func NewEventRepo() *EventRepo {
	return &EventRepo{
		es: make(map[string]domain.Event),
	}
}

type UserRepo struct {
	us map[string]args.SaveUser
}

func (r UserRepo) Save(ctx context.Context, su args.SaveUser) error {
	r.us[su.Login] = su
	return nil
}

func (r UserRepo) Login(ctx context.Context, l args.Login) (result.Login, error) {
	su, ok := r.us[l.Login]
	if !ok || su.PasswordHash != l.PasswordHash {
		return result.Login{}, errors.New("invalid cred")
	}
	return result.Login{UserID: 1}, nil
}

func (r UserRepo) HasUser(login string) bool { // test only
	_, ok := r.us[login]
	return ok
}

func NewUserRepo() *UserRepo {
	ur := &UserRepo{
		us: make(map[string]args.SaveUser),
	}
	return ur
}

func injectTestApp() fx.Option {
	return fx.Provide(
		NewConfig,

		NewUserRepo,
		fx.Annotate(func(r *UserRepo) *UserRepo {
			return r
		}, fx.As(new(repo.User))),

		NewEventRepo,
		fx.Annotate(func(r *EventRepo) *EventRepo {
			return r
		}, fx.As(new(repo.Event))),

		event.NewBus,
		event.NewSaveUserHandler,
		event.NewLoginHandler,

		command.NewBus,
		command.NewSaveUserHandler,
		command.NewLoginHandler,
		command.NewCalcAccrualHandler,

		handler.NewList,
		router.NewMux,
	)
}

func initTest(t *testing.T, r interface{}) {
	app := fxtest.New(t, injectTestApp(), fx.Invoke(r))
	defer app.RequireStop()
	app.RequireStart()
}
