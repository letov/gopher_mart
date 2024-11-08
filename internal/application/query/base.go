package query

import (
	"context"
	"errors"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"sync"
)

var (
	ErrNoHandler = errors.New("no handler for query")
)

type Name = string

type Query interface {
	GetName() Name
}

type Handler interface {
	Execute(q Query) (interface{}, error)
}

type Bus struct {
	sync.RWMutex
	handlers map[Name]Handler
	log      *zap.SugaredLogger
}

func NewBus(
	lc fx.Lifecycle,
	getOrdersHandler *GetOrdersHandler,
	log *zap.SugaredLogger,
) *Bus {
	b := &Bus{
		handlers: make(map[Name]Handler),
		log:      log,
	}

	b.Register(getOrdersHandler, GetOrdersName)

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

func (b *Bus) Execute(q Query) (interface{}, error) {
	b.RLock()
	h, ok := b.handlers[q.GetName()]
	b.RUnlock()
	if !ok {
		return nil, ErrNoHandler
	}
	res, err := h.Execute(q)
	if err != nil {
		b.log.Warnw("Error handling query", "name", q.GetName(), "err", err)
	}
	return res, err
}
