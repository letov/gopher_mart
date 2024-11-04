package service

import (
	"context"
	"errors"
	"gopher_mart/internal/application/dto/in"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/dto/response"
	"gopher_mart/internal/infrastructure/httpclient"
	"gopher_mart/internal/infrastructure/queue"
	"gopher_mart/internal/infrastructure/repo"
	"time"
)

var (
	ErrCantMapOrderStatus = errors.New("can't map order status")
)

type Accrual struct {
	client    httpclient.Client
	queue     queue.RequestAccrual
	orderRepo repo.Order
}

func (s Accrual) GetRequestAccrualHandler() queue.RequestAccrualHandler {
	return func(orderID string) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		oa, err := s.client.GetAccrual(ctx, orderID)
		if err != nil {
			return err
		}
		uo, err := mapper(oa)
		if err != nil {
			return err
		}
		return s.orderRepo.UpdateOrder(ctx, uo)
	}
}

func mapper(res response.Order) (in.UpdateOrder, error) {
	var os domain.OrderStatus

	switch res.Status {
	case response.RegisteredStatus:
		os = domain.NewStatus
	case response.InvalidStatus:
		os = domain.InvalidStatus
	case response.ProcessedStatus:
		os = domain.ProcessedStatus
	case response.ProcessingStatus:
		os = domain.ProcessingStatus
	default:
		return in.UpdateOrder{}, ErrCantMapOrderStatus
	}

	return in.UpdateOrder{
		OrderID: res.OrderID,
		Status:  os,
		Accrual: res.Accrual,
	}, nil
}

func (s Accrual) AddRequestToQueue(ctx context.Context, data in.RequestAccrual) error {
	return s.queue.Publish(ctx, data.OrderID)
}

func NewAccrual(
	client httpclient.Client,
	queue queue.RequestAccrual,
	orderRepo repo.Order,
) *Accrual {
	a := &Accrual{
		client,
		queue,
		orderRepo,
	}

	_ = queue.RegisterHandler(a.GetRequestAccrualHandler())

	return a
}
