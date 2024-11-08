package in

import "gopher_mart/internal/domain"

type SaveOperation struct {
	OrderID string
	UserID  int64
	Status  domain.OperationStatus
	Sum     int64
}
