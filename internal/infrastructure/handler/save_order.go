package handler

import (
	"context"
	"encoding/json"
	"errors"
	"gopher_mart/internal/application/command"
	"gopher_mart/internal/infrastructure/dto/request"
	"gopher_mart/internal/utils"
	"io"
	"net/http"
	"time"
)

const SaveOrderName string = "SaveOrderName"

func NewSaveOrderHandler(cb *command.Bus) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		data, err := io.ReadAll(req.Body)
		_ = req.Body.Close()
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		var dto request.SaveOrder
		err = json.Unmarshal(data, &dto)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		userId, err := utils.GetUserIdFromToken(req.Context())
		if err != nil {
			http.Error(res, err.Error(), http.StatusForbidden)
			return
		}

		cmd := command.SaveOrder{
			Ctx:     ctx,
			Request: dto,
			UserID:  userId,
		}
		_, err = cb.Execute(cmd)
		if err != nil {
			switch {
			case errors.Is(err, command.ErrInvalidOrderId):
				http.Error(res, err.Error(), http.StatusUnprocessableEntity)
			case errors.Is(err, command.ErrRequestByThisUserAlreadyExists):
				return
			case errors.Is(err, command.ErrRequestByAnotherUserAlreadyExists):
				http.Error(res, err.Error(), http.StatusConflict)
			default:
				http.Error(res, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		res.WriteHeader(http.StatusAccepted)
	}
}
