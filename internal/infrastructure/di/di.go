package di

import (
	"go.uber.org/fx"
	"gopher_mart/internal/infrastructure/config"
)

func InjectApp() fx.Option {
	return fx.Provide(
		config.NewConfig,
	)
}
