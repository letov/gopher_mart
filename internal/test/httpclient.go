package test

import (
	"context"
	"gopher_mart/internal/infrastructure/dto/response"
)

type HttpClient struct {
	res response.Order
}

func (h *HttpClient) SetResponse(res response.Order) {
	h.res = res
}

func (h *HttpClient) GetAccrual(ctx context.Context, orderId string) (response.Order, error) {
	return h.res, nil
}

func NewHttpClient() *HttpClient {
	return &HttpClient{}
}
