package event

import (
	"context"
	"gopher_mart/internal/application/dto"
	"gopher_mart/internal/domain"
)

const CreateUserEventName Name = "CreateUserEventName"

type CreateUserEvent struct {
	Ctx  context.Context
	User dto.CreateUser
}

func (e CreateUserEvent) GetName() Name {
	return CreateUserEventName
}

type CreateUserEventHandler struct {
	repo domain.UserRepository
}

func (h CreateUserEventHandler) Handle(e Event) {
	event := e.(CreateUserEvent)
	u := event.User
	_ = h.repo.Save(event.Ctx, domain.User{
		Login:        u.Login,
		PasswordHash: u.PasswordHash,
	})
}

func NewCreateUserEventHandler(repo domain.UserRepository) *CreateUserEventHandler {
	return &CreateUserEventHandler{
		repo,
	}
}
