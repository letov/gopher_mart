package event

import (
	"context"
	"errors"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/repo"
	"sync"
)

var (
	ErrNoHandler = errors.New("no handler for event")
)

type Event interface {
	GetAction() domain.Action
}

type Handler interface {
	Handle(e Event) error
}

type Bus struct {
	log *zap.SugaredLogger
	sync.RWMutex
	handlers map[domain.Action][]Handler
}

func NewBus(
	lc fx.Lifecycle,
	log *zap.SugaredLogger,
	saveUserHandler *SaveUserHandler,
	loginHandler *LoginHandler,
	saveOrderHandler *SaveOrderHandler,
	updateOrderHandler *UpdateOrderHandler,
	saveOperationHandler *SaveOperationHandler,
) *Bus {
	b := &Bus{
		handlers: make(map[domain.Action][]Handler),
		log:      log,
	}

	b.Subscribe(saveUserHandler, domain.SaveUserAction)
	b.Subscribe(loginHandler, domain.LoginAction)
	b.Subscribe(saveOrderHandler, domain.SaveOrderAction)
	b.Subscribe(updateOrderHandler, domain.UpdateOrderAction)
	b.Subscribe(saveOperationHandler, domain.SaveOperationAction)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			// TODO: stop active events
			return nil
		},
	})

	return b
}

func (b *Bus) Subscribe(h Handler, a domain.Action) {
	b.Lock()
	hs, ok := b.handlers[a]
	if ok {
		b.handlers[a] = append(hs, h)
	} else {
		b.handlers[a] = []Handler{h}
	}
	b.Unlock()
}

func (b *Bus) Publish(e Event) error {
	b.RLock()
	hs, ok := b.handlers[e.GetAction()]
	b.RUnlock()
	var wg sync.WaitGroup
	if !ok {
		return ErrNoHandler
	}
	for _, h := range hs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := h.Handle(e)
			if err != nil {
				b.log.Warnw("Error handling event", "name", e.GetAction(), "err", err)
			}
		}()
	}
	wg.Wait()
	return nil
}

type BaseHandler struct {
	eventRepo repo.Event
}

func (h *BaseHandler) Save(ctx context.Context, e domain.Event) error {
	return h.eventRepo.Save(ctx, e)
}

func NewBaseHandler(eventRepo repo.Event) *BaseHandler {
	return &BaseHandler{
		eventRepo: eventRepo,
	}
}
