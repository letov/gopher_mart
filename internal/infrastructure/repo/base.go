package repo

import (
	"context"
	"gopher_mart/internal/application/dto/in"
	"gopher_mart/internal/domain"
)

type User interface {
	Get(ctx context.Context, login string) (domain.User, error)
	Save(ctx context.Context, su in.SaveUser) error
}

type Order interface {
	Get(ctx context.Context, orderId string) (domain.Order, error)
	GetByUserId(ctx context.Context, userId int64) ([]domain.Order, error)
	Save(ctx context.Context, ra in.RequestAccrual) error
	UpdateOrder(ctx context.Context, uo in.UpdateOrder) error
}

type Event interface {
	GetLast(ctx context.Context, rootID string, action domain.Action) (domain.Event, error)
	Save(ctx context.Context, e domain.Event) error
}

type Operation interface {
}
