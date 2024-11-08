package event

import (
	"context"
	"gopher_mart/internal/application/dto/in"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/repo"
)

type UpdateOrder struct {
	Ctx  context.Context
	Data in.UpdateOrder
}

func (e UpdateOrder) GetAction() domain.Action {
	return domain.UpdateOrderAction
}

type UpdateOrderHandler struct {
	bh   *BaseHandler
	repo repo.Order
}

func (h UpdateOrderHandler) Handle(e Event) error {
	event := e.(UpdateOrder)
	err := h.bh.Save(event.Ctx, domain.Event{
		RootID:  event.Data.OrderID,
		Action:  domain.UpdateOrderAction,
		Payload: event.Data,
	})
	if err != nil {
		return err
	}

	return h.repo.Update(event.Ctx, event.Data)
}

func NewUpdateOrderHandler(
	bh *BaseHandler,
	repo repo.Order,
) *UpdateOrderHandler {
	return &UpdateOrderHandler{
		bh:   bh,
		repo: repo,
	}
}
