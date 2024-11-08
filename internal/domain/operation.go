package domain

import "time"

type OperationStatus string

const (
	AddedStatus    OperationStatus = "ADDED"
	DeductedStatus OperationStatus = "DEDUCTED"
)

type Operation struct {
	ID        int64
	OrderID   string
	UserID    int64
	Status    OperationStatus
	Sum       int64
	CreatedAt time.Time
}
