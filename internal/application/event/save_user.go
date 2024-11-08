package event

import (
	"context"
	"gopher_mart/internal/application/dto/in"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/repo"
)

type SaveUser struct {
	Ctx  context.Context
	Data in.SaveUser
}

func (e SaveUser) GetAction() domain.Action {
	return domain.SaveUserAction
}

type SaveUserHandler struct {
	bh   *BaseHandler
	repo repo.User
}

func (h SaveUserHandler) Handle(e Event) error {
	event := e.(SaveUser)
	err := h.bh.Save(event.Ctx, domain.Event{
		RootID:  event.Data.Login,
		Action:  domain.SaveUserAction,
		Payload: event.Data,
	})
	if err != nil {
		return err
	}

	return h.repo.Save(event.Ctx, event.Data)
}

func NewSaveUserHandler(bh *BaseHandler, repo repo.User) *SaveUserHandler {
	return &SaveUserHandler{
		bh:   bh,
		repo: repo,
	}
}
