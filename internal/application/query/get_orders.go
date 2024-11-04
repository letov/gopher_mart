package query

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"gopher_mart/internal/application/command"
	"gopher_mart/internal/infrastructure/config"
	"gopher_mart/internal/infrastructure/repo"
)

var (
	ErrHasNoOrders = errors.New("user has no orders")
)

const GetOrdersName command.Name = "GetOrdersName"

type GetOrders struct {
	Ctx    context.Context
	UserID int64
}

func (c GetOrders) GetName() command.Name {
	return GetOrdersName
}

type GetOrdersHandler struct {
	config    *config.Config
	orderRepo repo.Order
}

func (h GetOrdersHandler) Execute(q Query) (interface{}, error) {
	query := q.(GetOrders)

	os, err := h.orderRepo.GetByUserId(query.Ctx, query.UserID)

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrHasNoOrders
	}
	if err != nil {
		return nil, err
	}

	return os, nil
}

func NewGetOrdersHandler(
	config *config.Config,
	orderRepo repo.Order,
) *GetOrdersHandler {
	return &GetOrdersHandler{
		config:    config,
		orderRepo: orderRepo,
	}
}
