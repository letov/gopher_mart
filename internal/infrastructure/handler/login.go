package handler

import (
	"context"
	"encoding/json"
	"gopher_mart/internal/application/command"
	"gopher_mart/internal/infrastructure/dto/request"
	"gopher_mart/internal/infrastructure/dto/response"
	"io"
	"net/http"
	"time"
)

const LoginName string = "LoginName"

func NewLoginHandler(cb *command.Bus) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		data, err := io.ReadAll(req.Body)
		_ = req.Body.Close()
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		var dto request.Login
		err = json.Unmarshal(data, &dto)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		cmd := command.Login{
			Ctx:     ctx,
			Request: dto,
		}

		result, err := cb.Execute(cmd)
		if err != nil {
			http.Error(res, err.Error(), http.StatusForbidden)
			return
		}

		body, err := json.Marshal(result.(response.Login))
		res.Header().Set("Content-Type", "application/json")
		if _, err = res.Write(body); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
