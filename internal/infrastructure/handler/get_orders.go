package handler

import (
	"context"
	"encoding/json"
	"errors"
	"gopher_mart/internal/application/query"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/dto/response"
	"gopher_mart/internal/utils"
	"net/http"
	"time"
)

const GetOrdersName string = "GetOrdersName"

func mapper(dos []domain.Order) []response.Order {
	var res []response.Order
	for _, do := range dos {
		ro := response.Order{
			OrderID: do.OrderID,
			Status:  do.Status,
		}
		if do.Accrual > 0 {
			ro.Accrual = do.Accrual
		}
		ro.CreatedAt = do.CreatedAt.Format(time.RFC3339)
		res = append(res, ro)
	}
	return res
}

func NewGetOrdersHandler(qb *query.Bus) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		userId, err := utils.GetUserIdFromToken(req.Context())
		if err != nil {
			http.Error(res, err.Error(), http.StatusForbidden)
			return
		}

		q := query.GetOrders{
			Ctx:    ctx,
			UserID: userId,
		}
		os, err := qb.Execute(q)

		if err != nil {
			switch {
			case errors.Is(err, query.ErrHasNoOrders):
				http.Error(res, err.Error(), http.StatusNoContent)
			default:
				http.Error(res, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		data := mapper(os.([]domain.Order))
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
