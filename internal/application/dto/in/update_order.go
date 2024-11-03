package in

import "gopher_mart/internal/domain"

type UpdateOrder struct {
	OrderID string
	Status  domain.OrderStatus
	Accrual int64
}
