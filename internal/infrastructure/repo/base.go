package repo

import (
	"context"
	"gopher_mart/internal/application/dto/args"
	"gopher_mart/internal/application/dto/result"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/dto/response"
	"time"
)

type User interface {
	Save(ctx context.Context, su args.SaveUser) error
	Login(ctx context.Context, l args.Login) (result.Login, error)
}

type Order interface {
	Save(ctx context.Context, su args.CalcAccrual) error
	UpdateOrder(ctx context.Context, oa response.OrderAccrual) error
}

type Event interface {
	Save(ctx context.Context, e domain.Event) error
	HasEvent(ctx context.Context, rootID string, action domain.Action) bool
	HasEventWithDuration(ctx context.Context, rootID string, action domain.Action, duration time.Duration) bool
}

type Operation interface {
}
