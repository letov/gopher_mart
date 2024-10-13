package event

import (
	"context"
	"gopher_mart/internal/application/dto"
	"gopher_mart/internal/domain"
)

const SaveUserName Name = "SaveUserName"

type SaveUser struct {
	Ctx  context.Context
	User dto.SaveUser
}

func (e SaveUser) GetName() Name {
	return SaveUserName
}

type SaveUserHandler struct {
	repo domain.UserRepository
}

func (h SaveUserHandler) Handle(e Event) {
	event := e.(SaveUser)
	u := event.User
	_ = h.repo.Save(event.Ctx, domain.User{
		Login:        u.Login,
		PasswordHash: u.PasswordHash,
	})
}

func NewSaveUserHandler(repo domain.UserRepository) *SaveUserHandler {
	return &SaveUserHandler{
		repo,
	}
}
