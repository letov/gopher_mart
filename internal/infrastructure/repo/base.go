package repo

import (
	"context"
	"gopher_mart/internal/application/dto/args"
	"gopher_mart/internal/domain"
)

type User interface {
	Save(ctx context.Context, su args.SaveUser) error
	Login(ctx context.Context, l args.Login) bool
}

type Event interface {
	Save(ctx context.Context, e domain.Event) error
	HasEvent(ctx context.Context, rootID string, action domain.Action) bool
}
