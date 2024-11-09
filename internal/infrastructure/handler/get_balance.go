package handler

import (
	"context"
	"encoding/json"
	"gopher_mart/internal/application/query"
	"gopher_mart/internal/infrastructure/dto/response"
	"gopher_mart/internal/utils"
	"net/http"
	"time"
)

const GetBalanceName string = "GetBalanceName"

func NewGetBalanceHandler(qb *query.Bus) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userId, err := utils.GetUserIdFromToken(req.Context())
		if err != nil {
			http.Error(res, err.Error(), http.StatusForbidden)
			return
		}

		q := query.GetBalance{
			Ctx:    ctx,
			UserID: userId,
		}
		b, err := qb.Execute(q)

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		ob := b.(response.Balance)
		data := response.Balance{
			Current:   ob.Current,
			Withdrawn: ob.Withdrawn,
		}
		body, err := json.Marshal(data)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		if _, err = res.Write(body); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
