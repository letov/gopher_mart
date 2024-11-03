package repo

import (
	"context"
	"gopher_mart/internal/application/dto/in"
	"gopher_mart/internal/application/dto/out"
	"gopher_mart/internal/domain"
	"time"
)

type User interface {
	Save(ctx context.Context, su in.SaveUser) error
	Login(ctx context.Context, l in.Login) (out.Login, error)
}

type Order interface {
	Get(ctx context.Context, orderId string) (domain.Order, error)
	Save(ctx context.Context, ra in.RequestAccrual) error
	UpdateOrder(ctx context.Context, uo in.UpdateOrder) error
}

type Event interface {
	Save(ctx context.Context, e domain.Event) error
	HasEvent(ctx context.Context, rootID string, action domain.Action, durationSec time.Duration) (bool, error)
}

type Operation interface {
}
