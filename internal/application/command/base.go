package command

import (
	"context"
	"errors"
	"go.uber.org/fx"
	"go.uber.org/zap"
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
	log      *zap.SugaredLogger
}

func NewBus(
	lc fx.Lifecycle,
	saveUserHandler *SaveUserHandler,
	loginHandler *LoginHandler,
	requestAccrualHandler *SaveOrderHandler,
	log *zap.SugaredLogger,
) *Bus {
	b := &Bus{
		handlers: make(map[Name]Handler),
		log:      log,
	}

	b.Register(saveUserHandler, SaveUserName)
	b.Register(loginHandler, LoginName)
	b.Register(requestAccrualHandler, SaveOrderName)

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
	res, err := h.Execute(c)
	if err != nil {
		b.log.Warn(err)
	}
	return res, err
}
