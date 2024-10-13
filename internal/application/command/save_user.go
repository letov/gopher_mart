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

const SaveUserName Name = "SaveUser"

type SaveUser struct {
	Ctx      context.Context
	Login    string
	Password string
}

func (c SaveUser) GetName() Name {
	return SaveUserName
}

type SaveUserHandler struct {
	repo     domain.EventRepository
	config   *config.Config
	eventBus *event.Bus
}

func (h SaveUserHandler) Execute(c Command) (interface{}, error) {
	salt := h.config.Salt
	cmd := c.(SaveUser)

	if h.repo.HasEvent(cmd.Ctx, cmd.Login, domain.SaveUser) {
		return nil, ErrUserExist
	}

	u := dto.SaveUser{
		Login:        cmd.Login,
		PasswordHash: utils.GetHash(cmd.Password, salt),
	}

	if err := h.repo.Save(cmd.Ctx, domain.Event{
		RootID:  u.Login,
		Action:  domain.SaveUser,
		Payload: u,
	}); err != nil {
		return nil, err
	}

	err := h.eventBus.Publish(event.SaveUser{
		Ctx:  cmd.Ctx,
		User: u,
	})

	return u, err
}

func NewSaveUserHandler(
	repo domain.EventRepository,
	config *config.Config,
	eventBus *event.Bus,
) *SaveUserHandler {
	return &SaveUserHandler{
		repo,
		config,
		eventBus,
	}
}
