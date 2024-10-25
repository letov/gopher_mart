package handler

import (
	"context"
	"encoding/json"
	"github.com/go-chi/jwtauth"
	"gopher_mart/internal/application/command"
	"gopher_mart/internal/infrastructure/dto/request"
	"io"
	"net/http"
	"time"
)

const CalcAccrualName string = "CalcAccrualName"

func NewCalcAccrualHandler(cb *command.Bus) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		data, err := io.ReadAll(req.Body)
		_ = req.Body.Close()
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		var dto request.CalcAccrual
		err = json.Unmarshal(data, &dto)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		_, claims, err := jwtauth.FromContext(req.Context())
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		cmd := command.CalcAccrual{
			Ctx:     ctx,
			Request: dto,
			UserID:  claims["user_id"].(int64),
		}
		_, err = cb.Execute(cmd)
		if err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}
	}
}
