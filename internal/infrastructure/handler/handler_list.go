package handler

import (
	"errors"
	"gopher_mart/internal/application/command"
	"gopher_mart/internal/application/query"
	"net/http"
)

var (
	ErrNoHandler = errors.New("no handler for route")
)

type List struct {
	handlers map[string]http.HandlerFunc
}

func (l List) Get(n string) http.HandlerFunc {
	h, ok := l.handlers[n]
	if !ok {
		panic(ErrNoHandler)
	}
	return h
}

func NewList(
	cb *command.Bus,
	qb *query.Bus,
) *List {
	handlers := make(map[string]http.HandlerFunc)

	handlers[SaveUserName] = NewSaveUserHandler(cb)
	handlers[LoginName] = NewLoginHandler(cb)
	handlers[SaveOrderName] = NewSaveOrderHandler(cb)
	handlers[GetOrdersName] = NewGetOrdersHandler(qb)
	handlers[GetBalanceName] = NewGetBalanceHandler(qb)

	return &List{
		handlers,
	}
}
