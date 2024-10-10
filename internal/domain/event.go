package domain

import (
	"context"
)

type Actions string

const (
	CreateUser Actions = "CreateUser"
)

type Event struct {
	RootID  string
	Action  Actions
	Payload interface{}
}

type EventRepository interface {
	HasEvent(ctx context.Context, rootID string, action Actions) bool
	Save(ctx context.Context, e Event) error
}
