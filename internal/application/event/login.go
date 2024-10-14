package event

import (
	"context"
	"gopher_mart/internal/application/dto/args"
)

const LoginName Name = "LoginNameName"

type Login struct {
	Ctx  context.Context
	Data args.Login
}

func (e Login) GetName() Name {
	return LoginName
}

type LoginHandler struct{}

func (h LoginHandler) Handle(e Event) {}

func NewLoginHandler() *LoginHandler {
	return &LoginHandler{}
}
