package httpclient

import (
	"context"
	"gopher_mart/internal/infrastructure/dto/response"
)

type OrderAccrual interface {
	GetAccrual(ctx context.Context, orderId int64) (response.OrderAccrual, error)
}
