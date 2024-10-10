package command

import (
	"context"
	"errors"
	"gopher_mart/internal/application/dto"
	"gopher_mart/internal/application/event"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/config"
	"gopher_mart/internal/utils"
)

var (
	ErrUserExist = errors.New("user already exists")
)

const CreateUserCommandName Name = "CreateUserCommand"

type CreateUserCommand struct {
	Ctx      context.Context
	Login    string
	Password string
}

func (c CreateUserCommand) GetName() Name {
	return CreateUserCommandName
}

type CreateUserCommandHandler struct {
	repo     domain.EventRepository
	config   *config.Config
	eventBus *event.Bus
}

func (h CreateUserCommandHandler) Execute(c Command) (interface{}, error) {
	salt := h.config.Salt
	cmd := c.(CreateUserCommand)

	if h.repo.HasEvent(cmd.Ctx, cmd.Login, domain.CreateUser) {
		return nil, ErrUserExist
	}

	u := dto.CreateUser{
		Login:        cmd.Login,
		PasswordHash: utils.GetHash(cmd.Password, salt),
	}

	if err := h.repo.Save(cmd.Ctx, domain.Event{
		RootID:  u.Login,
		Action:  domain.CreateUser,
		Payload: u,
	}); err != nil {
		return nil, err
	}

	err := h.eventBus.Publish(event.CreateUserEvent{
		Ctx:  cmd.Ctx,
		User: u,
	})

	return u, err
}

func NewCreateUserCommandHandler(
	repo domain.EventRepository,
	config *config.Config,
	eventBus *event.Bus,
) *CreateUserCommandHandler {
	return &CreateUserCommandHandler{
		repo,
		config,
		eventBus,
	}
}
