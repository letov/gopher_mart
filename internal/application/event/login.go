package event

import (
	"context"
	"gopher_mart/internal/application/dto/in"
	"gopher_mart/internal/domain"
)

type Login struct {
	Ctx  context.Context
	Data in.Login
}

func (e Login) GetAction() domain.Action {
	return domain.LoginAction
}

type LoginHandler struct {
	bh *BaseHandler
}

func (h LoginHandler) Handle(e Event) error {
	event := e.(Login)
	return h.bh.Save(event.Ctx, domain.Event{
		RootID:  event.Data.Login,
		Action:  domain.LoginAction,
		Payload: event.Data,
	})
}

func NewLoginHandler(bh *BaseHandler) *LoginHandler {
	return &LoginHandler{
		bh,
	}
}
