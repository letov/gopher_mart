package command

import (
	"context"
	"errors"
	"gopher_mart/internal/application/dto/args"
	"gopher_mart/internal/application/event"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/config"
	"gopher_mart/internal/infrastructure/dto/request"
	"gopher_mart/internal/infrastructure/repo"
	"gopher_mart/internal/utils"
)

var (
	ErrUserExist = errors.New("user already exists")
)

const SaveUserName Name = "SaveUserAction"

type SaveUser struct {
	Ctx     context.Context
	Request request.SaveUser
}

func (c SaveUser) GetName() Name {
	return SaveUserName
}

type SaveUserHandler struct {
	config    *config.Config
	eventRepo repo.Event
	eventBus  *event.Bus
}

func (h SaveUserHandler) Execute(c Command) (interface{}, error) {
	salt := h.config.Salt
	cmd := c.(SaveUser)
	hash := utils.GetHash(cmd.Request.Password, salt)

	he, err := h.eventRepo.HasEvent(cmd.Ctx, cmd.Request.Login, domain.SaveUserAction, 0)
	if err != nil {
		return nil, err
	} else if he {
		return nil, ErrUserExist
	}

	data := args.SaveUser{
		Login:        cmd.Request.Login,
		PasswordHash: hash,
		Name:         cmd.Request.Name,
	}

	if err := h.eventRepo.Save(cmd.Ctx, domain.Event{
		RootID:  data.Login,
		Action:  domain.SaveUserAction,
		Payload: data,
	}); err != nil {
		return nil, err
	}

	err = h.eventBus.Publish(event.SaveUser{
		Ctx:  cmd.Ctx,
		Data: data,
	})

	return data, err
}

func NewSaveUserHandler(
	config *config.Config,
	eventRepo repo.Event,
	eventBus *event.Bus,
) *SaveUserHandler {
	return &SaveUserHandler{
		config:    config,
		eventRepo: eventRepo,
		eventBus:  eventBus,
	}
}
