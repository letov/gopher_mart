package handler

import (
	"context"
	"encoding/json"
	"errors"
	"gopher_mart/internal/application/command"
	"gopher_mart/internal/infrastructure/dto/request"
	"io"
	"net/http"
	"time"
)

const SaveUserName string = "SaveUserName"

func NewSaveUserHandler(cb *command.Bus) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		data, err := io.ReadAll(req.Body)
		_ = req.Body.Close()
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		var dto request.SaveUser
		err = json.Unmarshal(data, &dto)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		cmd := command.SaveUser{
			Ctx:     ctx,
			Request: dto,
		}
		_, err = cb.Execute(cmd)
		if err != nil {
			switch {
			case errors.Is(err, command.ErrUserExists):
				http.Error(res, err.Error(), http.StatusConflict)
			default:
				http.Error(res, err.Error(), http.StatusInternalServerError)
			}
			return
		}
	}
}
