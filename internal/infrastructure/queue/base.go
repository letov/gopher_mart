package queue

import "context"

type CalcAccrualHandler = func(orderId int64) error

type CalcAccrual interface {
	Publish(ctx context.Context, orderId int64) error
	RegisterHandler(h CalcAccrualHandler) error
}
