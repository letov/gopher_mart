package command

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"gopher_mart/internal/application/dto/in"
	"gopher_mart/internal/application/event"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/config"
	"gopher_mart/internal/infrastructure/dto/request"
	"gopher_mart/internal/infrastructure/repo"
	"gopher_mart/internal/utils"
)

var (
	ErrInvalidOrderId                    = errors.New("invalid order id")
	ErrRequestByThisUserAlreadyExists    = errors.New("request by this user already exists")
	ErrRequestByAnotherUserAlreadyExists = errors.New("request by another user already exists")
)

const SaveOrderName Name = "SaveOrderName"

type SaveOrder struct {
	Ctx     context.Context
	Request request.SaveOrder
	UserID  int64
}

func (c SaveOrder) GetName() Name {
	return SaveOrderName
}

type SaveOrderHandler struct {
	config    *config.Config
	eventRepo repo.Event
	eventBus  *event.Bus
}

func (h SaveOrderHandler) Execute(c Command) (interface{}, error) {
	cmd := c.(SaveOrder)
	if !utils.IsValidOrderId(cmd.Request.OrderID) {
		return nil, ErrInvalidOrderId
	}

	e, err := h.eventRepo.GetLast(
		cmd.Ctx,
		cmd.Request.OrderID,
		domain.SaveOrderAction,
	)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	if err == nil {
		u := e.Payload.(map[string]interface{})["UserID"].(float64)
		if int64(u) == cmd.UserID {
			return nil, ErrRequestByThisUserAlreadyExists
		} else {
			return nil, ErrRequestByAnotherUserAlreadyExists
		}
	}

	data := in.RequestAccrual{
		OrderID: cmd.Request.OrderID,
		UserID:  cmd.UserID,
	}

	if err := h.eventRepo.Save(cmd.Ctx, domain.Event{
		RootID:  cmd.Request.OrderID,
		Action:  domain.SaveOrderAction,
		Payload: data,
	}); err != nil {
		return nil, err
	}

	err = h.eventBus.Publish(event.SaveOrder{
		Ctx:  cmd.Ctx,
		Data: data,
	})

	return data, err
}

func NewSaveOrderHandler(
	config *config.Config,
	eventRepo repo.Event,
	eventBus *event.Bus,
) *SaveOrderHandler {
	return &SaveOrderHandler{
		config:    config,
		eventRepo: eventRepo,
		eventBus:  eventBus,
	}
}
