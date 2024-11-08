package service

import (
	"context"
	"errors"
	"gopher_mart/internal/application/dto/in"
	"gopher_mart/internal/application/event"
	"gopher_mart/internal/domain"
	"gopher_mart/internal/infrastructure/dto/response"
	"gopher_mart/internal/infrastructure/httpclient"
	"gopher_mart/internal/infrastructure/queue"
	"gopher_mart/internal/infrastructure/repo"
	"time"
)

var (
	ErrCantMapOrderStatus = errors.New("can't map order status")
	ErrWaitForFinalStatus = errors.New("wait for final")
)

type Accrual struct {
	client    httpclient.Client
	queue     queue.RequestAccrual
	eventBus  *event.Bus
	orderRepo repo.Order
}

// GetRequestAccrualHandler получаем из очереди заказ, у которого нужно запросить баллы
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

		err = s.eventBus.Publish(event.UpdateOrder{
			Ctx:  ctx,
			Data: uo,
		})
		if err != nil {
			return err
		}

		if !isFinalAccrualStatus(oa.Status) {
			return ErrWaitForFinalStatus
		} else {
			return s.saveOperation(ctx, oa)
		}
	}
}

func (s Accrual) saveOperation(ctx context.Context, ao response.AccrualOrder) error {
	o, err := s.orderRepo.Get(ctx, ao.OrderID)
	if err != nil {
		return err
	}

	data := in.SaveOperation{
		OrderID: ao.OrderID,
		UserID:  o.UserID,
		Status:  domain.AddedStatus,
		Sum:     ao.Accrual,
	}
	return s.eventBus.Publish(event.SaveOperation{
		Ctx:  ctx,
		Data: data,
	})
}

func isFinalAccrualStatus(status response.Status) bool {
	return status == response.ProcessedStatus || status == response.InvalidStatus
}

func mapper(res response.AccrualOrder) (in.UpdateOrder, error) {
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

func NewAccrual(
	client httpclient.Client,
	queue queue.RequestAccrual,
	eventBus *event.Bus,
	orderRepo repo.Order,
) *Accrual {
	a := &Accrual{
		client,
		queue,
		eventBus,
		orderRepo,
	}

	_ = queue.RegisterHandler(a.GetRequestAccrualHandler())

	return a
}
