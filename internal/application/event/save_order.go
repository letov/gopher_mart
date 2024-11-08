package event

import (
	"context"
	"gopher_mart/internal/application/dto/in"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/queue"
	"gopher_mart/internal/infrastructure/repo"
)

type SaveOrder struct {
	Ctx  context.Context
	Data in.RequestAccrual
}

func (e SaveOrder) GetAction() domain.Action {
	return domain.SaveOrderAction
}

type SaveOrderHandler struct {
	bh    *BaseHandler
	repo  repo.Order
	queue queue.RequestAccrual
}

func (h SaveOrderHandler) Handle(e Event) error {
	event := e.(SaveOrder)
	err := h.bh.Save(event.Ctx, domain.Event{
		RootID:  event.Data.OrderID,
		Action:  domain.SaveOrderAction,
		Payload: event.Data,
	})
	if err != nil {
		return err
	}

	err = h.repo.Save(event.Ctx, event.Data)
	if err != nil {
		return err
	}

	err = h.queue.Publish(event.Ctx, event.Data.OrderID)
	if err != nil {
		return err
	}

	return nil
}

func NewSaveOrderHandler(
	bh *BaseHandler,
	repo repo.Order,
	queue queue.RequestAccrual,
) *SaveOrderHandler {
	return &SaveOrderHandler{
		bh:    bh,
		repo:  repo,
		queue: queue,
	}
}
