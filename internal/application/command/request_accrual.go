package command

import (
	"context"
	"errors"
	"gopher_mart/internal/application/dto/in"
	"gopher_mart/internal/application/event"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/config"
	"gopher_mart/internal/infrastructure/dto/request"
	"gopher_mart/internal/infrastructure/repo"
	"gopher_mart/internal/utils"
)

var (
	ErrInvalidOrderId       = errors.New("invalid order id")
	ErrRequestAlreadyExists = errors.New("request already exists")
)

const RequestAccrualName Name = "RequestAccrualName"

type RequestAccrual struct {
	Ctx     context.Context
	Request request.RequestAccrual
	UserID  int64
}

func (c RequestAccrual) GetName() Name {
	return RequestAccrualName
}

type RequestAccrualHandler struct {
	config    *config.Config
	eventRepo repo.Event
	eventBus  *event.Bus
}

func (h RequestAccrualHandler) Execute(c Command) (interface{}, error) {
	cmd := c.(RequestAccrual)
	if !utils.IsValidOrderId(cmd.Request.OrderID) {
		return nil, ErrInvalidOrderId
	}

	he, err := h.eventRepo.HasEvent(
		cmd.Ctx,
		cmd.Request.OrderID,
		domain.RequestAccrualAction,
		0,
	)
	if err != nil {
		return nil, err
	} else if he {
		return nil, ErrRequestAlreadyExists
	}

	data := in.RequestAccrual{
		OrderID: cmd.Request.OrderID,
		UserID:  cmd.UserID,
	}

	if err := h.eventRepo.Save(cmd.Ctx, domain.Event{
		RootID:  cmd.Request.OrderID,
		Action:  domain.RequestAccrualAction,
		Payload: data,
	}); err != nil {
		return nil, err
	}

	err = h.eventBus.Publish(event.RequestAccrual{
		Ctx:  cmd.Ctx,
		Data: data,
	})

	return data, err
}

func NewRequestAccrualHandler(
	config *config.Config,
	eventRepo repo.Event,
	eventBus *event.Bus,
) *RequestAccrualHandler {
	return &RequestAccrualHandler{
		config:    config,
		eventRepo: eventRepo,
		eventBus:  eventBus,
	}
}
