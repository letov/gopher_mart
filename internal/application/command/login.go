package command

import (
	"context"
	"errors"
	"gopher_mart/internal/application/dto/args"
	"gopher_mart/internal/application/event"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/config"
	"gopher_mart/internal/infrastructure/dto/request"
	"gopher_mart/internal/infrastructure/dto/response"
	"gopher_mart/internal/infrastructure/repo"
	"gopher_mart/internal/utils"
)

var (
	ErrIncorrectLoginOrPassword = errors.New("incorrect login or password")
)

const LoginName Name = "LoginName"

type Login struct {
	Ctx     context.Context
	Request request.Login
}

func (c Login) GetName() Name {
	return LoginName
}

type LoginHandler struct {
	config    *config.Config
	userRepo  repo.User
	eventRepo repo.Event
	eventBus  *event.Bus
}

func (h LoginHandler) Execute(c Command) (interface{}, error) {
	salt := h.config.Salt
	cmd := c.(Login)
	hash := utils.GetHash(cmd.Request.Password, salt)

	data := args.Login{
		Login:        cmd.Request.Login,
		PasswordHash: hash,
	}

	login, err := h.userRepo.Login(cmd.Ctx, data)

	if err != nil {
		return nil, ErrIncorrectLoginOrPassword
	}

	if err := h.eventRepo.Save(cmd.Ctx, domain.Event{
		RootID:  data.Login,
		Action:  domain.LoginAction,
		Payload: data,
	}); err != nil {
		return nil, err
	}

	if err := h.eventBus.Publish(event.Login{
		Ctx:  cmd.Ctx,
		Data: data,
	}); err != nil {
		return nil, err
	}

	token, err := utils.GetJwt(h.config.JwtKey, login)
	if err != nil {
		return nil, err
	}

	result := response.Login{
		Token: token,
	}

	return result, err
}

func NewLoginHandler(
	config *config.Config,
	userRepo repo.User,
	eventRepo repo.Event,
	eventBus *event.Bus,
) *LoginHandler {
	return &LoginHandler{
		config:    config,
		userRepo:  userRepo,
		eventRepo: eventRepo,
		eventBus:  eventBus,
	}
}
