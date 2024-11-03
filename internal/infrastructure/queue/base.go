package queue

import "context"

type RequestAccrualHandler = func(orderID string) error

type RequestAccrual interface {
	Publish(ctx context.Context, orderID string) error
	RegisterHandler(h RequestAccrualHandler) error
}
