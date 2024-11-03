package test

import (
	"context"
	"gopher_mart/internal/infrastructure/dto/response"
)

type HttpClient struct {
	res response.OrderAccrual
}

func (h *HttpClient) SetResponse(res response.OrderAccrual) {
	h.res = res
}

func (h *HttpClient) GetAccrual(ctx context.Context, orderId string) (response.OrderAccrual, error) {
	return h.res, nil
}

func NewHttpClient() *HttpClient {
	return &HttpClient{}
}
