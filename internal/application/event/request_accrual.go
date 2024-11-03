package event

import (
	"context"
	"gopher_mart/internal/application/dto/in"
	"gopher_mart/internal/application/service"
	"gopher_mart/internal/infrastructure/repo"
)

const RequestAccrualName Name = "RequestAccrualName"

type RequestAccrual struct {
	Ctx  context.Context
	Data in.RequestAccrual
}

func (e RequestAccrual) GetName() Name {
	return RequestAccrualName
}

type RequestAccrualHandler struct {
	service *service.Accrual
	repo    repo.Order
}

func (h RequestAccrualHandler) Handle(e Event) {
	event := e.(RequestAccrual)
	_ = h.repo.Save(event.Ctx, event.Data)
	_ = h.service.AddRequestToQueue(event.Ctx, event.Data)
}

func NewRequestAccrualHandler(
	service *service.Accrual,
	repo repo.Order,
) *RequestAccrualHandler {
	return &RequestAccrualHandler{
		service,
		repo,
	}
}
