package handler

import (
	"errors"
	"gopher_mart/internal/application/command"
	"net/http"
)

var (
	ErrNoHandler = errors.New("no handler for command")
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

func NewList(cb *command.Bus) *List {
	handlers := make(map[string]http.HandlerFunc)

	handlers[SaveUserName] = NewSaveUserHandler(cb)

	return &List{
		handlers,
	}
}
