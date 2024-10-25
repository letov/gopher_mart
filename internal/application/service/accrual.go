package service

import (
	"context"
	"gopher_mart/internal/application/dto/args"
	"gopher_mart/internal/infrastructure/httpclient"
	"gopher_mart/internal/infrastructure/queue"
	"gopher_mart/internal/infrastructure/repo"
	"time"
)

type Accrual struct {
	client    httpclient.OrderAccrual
	queue     queue.CalcAccrual
	orderRepo repo.Order
}

func (s Accrual) GetCalcAccrualHandler() queue.CalcAccrualHandler {
	return func(orderId int64) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		oa, err := s.client.GetAccrual(ctx, orderId)
		if err != nil {
			return err
		}
		return s.orderRepo.UpdateOrder(ctx, oa)
	}
}

func (s Accrual) SaveOrder(ctx context.Context, data args.CalcAccrual) error {
	err := s.orderRepo.Save(ctx, data)
	if err != nil {
		return err
	}

	return s.addRequestToQueue(ctx, data.OrderID)
}

func (s Accrual) addRequestToQueue(ctx context.Context, orderID int64) error {
	return s.queue.Publish(ctx, orderID)
}

func NewAccrual(
	queue queue.CalcAccrual,
	orderRepo repo.Order,
) *Accrual {
	a := &Accrual{
		queue:     queue,
		orderRepo: orderRepo,
	}

	_ = queue.RegisterHandler(a.GetCalcAccrualHandler())

	return a
}
