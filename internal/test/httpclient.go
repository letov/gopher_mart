package test

import (
	"context"
	"gopher_mart/internal/infrastructure/dto/response"
)

type HttpClient struct {
	res response.AccrualOrder
}

func (h *HttpClient) SetResponse(res response.AccrualOrder) {
	h.res = res
}

func (h *HttpClient) GetAccrual(ctx context.Context, orderId string) (response.AccrualOrder, error) {
	return h.res, nil
}

func NewHttpClient() *HttpClient {
	return &HttpClient{}
}
