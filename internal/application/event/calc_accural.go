package event

import (
	"context"
	"gopher_mart/internal/application/dto/args"
	"gopher_mart/internal/application/service"
)

const CalcAccrualName Name = "CalcAccrualName"

type CalcAccrual struct {
	Ctx  context.Context
	Data args.CalcAccrual
}

func (e CalcAccrual) GetName() Name {
	return CalcAccrualName
}

type CalcAccrualHandler struct {
	service *service.Accrual
}

func (h CalcAccrualHandler) Handle(e Event) {
	event := e.(CalcAccrual)
	_ = h.service.SaveOrder(event.Ctx, event.Data)
}

func NewCalcAccrualHandler(
	service *service.Accrual,
) *CalcAccrualHandler {
	return &CalcAccrualHandler{
		service: service,
	}
}
