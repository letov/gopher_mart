package event

import (
	"context"
	"gopher_mart/internal/application/dto/args"
	"gopher_mart/internal/infrastructure/repo"
)

const SaveUserName Name = "SaveUserName"

type SaveUser struct {
	Ctx  context.Context
	Data args.SaveUser
}

func (e SaveUser) GetName() Name {
	return SaveUserName
}

type SaveUserHandler struct {
	repo repo.User
}

func (h SaveUserHandler) Handle(e Event) {
	event := e.(SaveUser)
	_ = h.repo.Save(event.Ctx, event.Data)
}

func NewSaveUserHandler(repo repo.User) *SaveUserHandler {
	return &SaveUserHandler{
		repo,
	}
}
