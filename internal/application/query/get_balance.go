package query

import (
	"context"
	"gopher_mart/internal/application/command"
	"gopher_mart/internal/infrastructure/config"
	"gopher_mart/internal/infrastructure/repo"
)

const GetBalanceName command.Name = "GetBalanceName"

type GetBalance struct {
	Ctx    context.Context
	UserID int64
}

func (c GetBalance) GetName() command.Name {
	return GetBalanceName
}

type GetBalanceHandler struct {
	config *config.Config
	repo   repo.Operation
}

func (h GetBalanceHandler) Execute(q Query) (interface{}, error) {
	query := q.(GetBalance)

	return h.repo.GetBalance(query.Ctx, query.UserID)
}

func NewGetBalanceHandler(
	config *config.Config,
	repo repo.Operation,
) *GetBalanceHandler {
	return &GetBalanceHandler{
		config: config,
		repo:   repo,
	}
}
