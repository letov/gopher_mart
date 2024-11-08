package command

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
	"gopher_mart/internal/application/dto/in"
	"gopher_mart/internal/application/event"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/config"
	"gopher_mart/internal/infrastructure/dto/request"
	"gopher_mart/internal/infrastructure/repo"
	"gopher_mart/internal/utils"
)

var (
	ErrUserExists = errors.New("user already exists")
)

const SaveUserName Name = "SaveUserName"

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

	_, err := h.eventRepo.GetLast(cmd.Ctx, cmd.Request.Login, domain.SaveUserAction)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}
	if err == nil {
		return nil, ErrUserExists
	}

	data := in.SaveUser{
		Login:        cmd.Request.Login,
		PasswordHash: hash,
		Name:         cmd.Request.Name,
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
