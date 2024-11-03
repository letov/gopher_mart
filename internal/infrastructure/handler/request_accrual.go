package handler

import (
	"context"
	"encoding/json"
	"gopher_mart/internal/application/command"
	"gopher_mart/internal/infrastructure/dto/request"
	"gopher_mart/internal/utils"
	"io"
	"net/http"
	"time"
)

const RequestAccrualName string = "RequestAccrualName"

func NewRequestAccrualHandler(cb *command.Bus) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		data, err := io.ReadAll(req.Body)
		_ = req.Body.Close()
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		var dto request.RequestAccrual
		err = json.Unmarshal(data, &dto)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		userId, err := utils.GetUserIdFromToken(req.Context())
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		cmd := command.RequestAccrual{
			Ctx:     ctx,
			Request: dto,
			UserID:  userId,
		}
		_, err = cb.Execute(cmd)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
