package event

import (
	"context"
	"errors"
	"go.uber.org/fx"
	"sync"
)

var (
	ErrNoHandler = errors.New("no handler for event")
)

type Name = string

type Event interface {
	GetName() Name
}

type Handler interface {
	Handle(e Event)
}

type Bus struct {
	sync.RWMutex
	handlers map[Name][]Handler
}

func NewBus(
	lc fx.Lifecycle,
	saveUserEventHandler *SaveUserHandler,
	loginHandler *LoginHandler,
	requestAccrualHandler *SaveOrderHandler,
) *Bus {
	b := &Bus{
		handlers: make(map[Name][]Handler),
	}

	b.Subscribe(saveUserEventHandler, SaveUserName)
	b.Subscribe(loginHandler, LoginName)
	b.Subscribe(requestAccrualHandler, SaveOrderName)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			// TODO: stop active events
			return nil
		},
	})

	return b
}

func (b *Bus) Subscribe(h Handler, en Name) {
	b.Lock()
	hs, ok := b.handlers[en]
	if ok {
		b.handlers[en] = append(hs, h)
	} else {
		b.handlers[en] = []Handler{h}
	}
	b.Unlock()
}

func (b *Bus) Publish(e Event) error {
	b.RLock()
	hs, ok := b.handlers[e.GetName()]
	b.RUnlock()
	var wg sync.WaitGroup
	if !ok {
		return ErrNoHandler
	}
	for _, h := range hs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			h.Handle(e)
		}()
	}
	wg.Wait()
	return nil
}
