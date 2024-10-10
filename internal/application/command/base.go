package command

import (
	"context"
	"errors"
	"go.uber.org/fx"
	"sync"
)

var (
	ErrNoHandler = errors.New("no handler for command")
)

type Name = string

type Command interface {
	GetName() Name
}

type Handler interface {
	Execute(c Command) (interface{}, error)
}

type Bus struct {
	sync.RWMutex
	handlers map[Name]Handler
}

func NewBus(
	lc fx.Lifecycle,
	createUserHandler *CreateUserCommandHandler,
) *Bus {
	b := &Bus{
		handlers: make(map[Name]Handler),
	}

	b.Register(createUserHandler, CreateUserCommandName)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			// TODO: stop active commands
			return nil
		},
	})

	return b
}

func (b *Bus) Register(h Handler, cn Name) {
	b.Lock()
	b.handlers[cn] = h
	b.Unlock()
}

func (b *Bus) Execute(c Command) (interface{}, error) {
	b.RLock()
	h, ok := b.handlers[c.GetName()]
	b.RUnlock()
	if !ok {
		return nil, ErrNoHandler
	}
	return h.Execute(c)
}
