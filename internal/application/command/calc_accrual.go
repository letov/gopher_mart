package command

import (
	"context"
	"errors"
	"gopher_mart/internal/application/dto/args"
	"gopher_mart/internal/application/event"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/config"
	"gopher_mart/internal/infrastructure/dto/request"
	"gopher_mart/internal/infrastructure/repo"
	"gopher_mart/internal/utils"
	"strconv"
	"time"
)

var (
	ErrInvalidOrderId = errors.New("invalid order id")
	ErrTooMoreTries   = errors.New("too more calc accrual tries")
)

const CalcAccrualName Name = "CalcAccrualName"

type CalcAccrual struct {
	Ctx     context.Context
	Request request.CalcAccrual
	UserID  int64
}

func (c CalcAccrual) GetName() Name {
	return CalcAccrualName
}

type CalcAccrualHandler struct {
	config    *config.Config
	eventRepo repo.Event
	eventBus  *event.Bus
}

func (h CalcAccrualHandler) Execute(c Command) (interface{}, error) {
	cmd := c.(CalcAccrual)
	if utils.Valid(cmd.Request.OrderID) {
		return nil, ErrInvalidOrderId
	}

	orderId := strconv.FormatInt(cmd.Request.OrderID, 10)

	he, err := h.eventRepo.HasEvent(
		cmd.Ctx,
		orderId,
		domain.CalcAccrualAction,
		300*time.Second,
	)
	if err != nil {
		return nil, err
	} else if he {
		return nil, ErrTooMoreTries
	}

	data := args.CalcAccrual{
		OrderID: cmd.Request.OrderID,
		UserID:  cmd.UserID,
	}

	if err := h.eventRepo.Save(cmd.Ctx, domain.Event{
		RootID:  orderId,
		Action:  domain.CalcAccrualAction,
		Payload: data,
	}); err != nil {
		return nil, err
	}

	err = h.eventBus.Publish(event.CalcAccrual{
		Ctx:  cmd.Ctx,
		Data: data,
	})

	return data, err
}

func NewCalcAccrualHandler(
	config *config.Config,
	eventRepo repo.Event,
	eventBus *event.Bus,
) *CalcAccrualHandler {
	return &CalcAccrualHandler{
		config:    config,
		eventRepo: eventRepo,
		eventBus:  eventBus,
	}
}
