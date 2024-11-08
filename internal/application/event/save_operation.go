package event

import (
	"context"
	"gopher_mart/internal/application/dto/in"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/repo"
)

type SaveOperation struct {
	Ctx  context.Context
	Data in.SaveOperation
}

func (e SaveOperation) GetAction() domain.Action {
	return domain.SaveOperationAction
}

type SaveOperationHandler struct {
	bh   *BaseHandler
	repo repo.Operation
}

func (h SaveOperationHandler) Handle(e Event) error {
	event := e.(SaveOperation)
	err := h.bh.Save(event.Ctx, domain.Event{
		RootID:  event.Data.OrderID,
		Action:  domain.SaveOperationAction,
		Payload: event.Data,
	})
	if err != nil {
		return err
	}

	return h.repo.Save(event.Ctx, event.Data)
}

func NewSaveOperationHandler(
	bh *BaseHandler,
	repo repo.Operation,
) *SaveOperationHandler {
	return &SaveOperationHandler{
		bh:   bh,
		repo: repo,
	}
}
