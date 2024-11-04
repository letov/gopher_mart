package event

import (
	"context"
	"gopher_mart/internal/application/dto/in"
	"gopher_mart/internal/application/service"
	"gopher_mart/internal/infrastructure/repo"
)

const SaveOrderName Name = "SaveOrderName"

type SaveOrder struct {
	Ctx  context.Context
	Data in.RequestAccrual
}

func (e SaveOrder) GetName() Name {
	return SaveOrderName
}

type SaveOrderHandler struct {
	service *service.Accrual
	repo    repo.Order
}

func (h SaveOrderHandler) Handle(e Event) {
	event := e.(SaveOrder)
	_ = h.repo.Save(event.Ctx, event.Data)
	_ = h.service.AddRequestToQueue(event.Ctx, event.Data)
}

func NewSaveOrderHandler(
	service *service.Accrual,
	repo repo.Order,
) *SaveOrderHandler {
	return &SaveOrderHandler{
		service,
		repo,
	}
}
