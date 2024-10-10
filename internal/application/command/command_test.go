package command

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gopher_mart/internal/application/event"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/test"
	"testing"
	"time"
)

func injectTestApp() fx.Option {
	options := fx.Provide(
		test.NewConfig,

		test.NewUserRepo,
		fx.Annotate(func(r *test.UserRepo) *test.UserRepo {
			return r
		}, fx.As(new(domain.UserRepository))),

		test.NewEventRepo,
		fx.Annotate(func(r *test.EventRepo) *test.EventRepo {
			return r
		}, fx.As(new(domain.EventRepository))),

		event.NewCreateUserEventHandler,
		event.NewBus,

		NewBus,
		NewCreateUserCommandHandler,
	)

	err := fx.ValidateApp(options)
	println(err)
	return options
}

func inttest(t *testing.T, r interface{}) {
	app := fxtest.New(t, injectTestApp(), fx.Invoke(r))
	defer app.RequireStop()
	app.RequireStart()
}

func Test_CreateUserCommand(t *testing.T) {
	type args struct {
		login    string
		password string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "positive test #1",
			args: args{
				login:    "login",
				password: "password",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inttest(t, func(cb *Bus, ur *test.UserRepo, er *test.EventRepo) {
				cmd := CreateUserCommand{
					Ctx:      context.Background(),
					Login:    tt.args.login,
					Password: tt.args.password,
				}
				_, err := cb.Execute(cmd)
				assert.Equal(t, err, nil)

				assert.True(t, er.HasRootId(tt.args.login))
				time.Sleep(time.Millisecond * 200) // UserRepo update async by event
				assert.True(t, ur.HasLogin(tt.args.login))
			})
		})
	}
}
