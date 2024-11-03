package httpclient

import (
	"context"
	"gopher_mart/internal/infrastructure/dto/response"
)

type Client interface {
	GetAccrual(ctx context.Context, orderID string) (response.OrderAccrual, error)
}
