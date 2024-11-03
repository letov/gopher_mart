package domain

import "time"

type OrderStatus string

const (
	NewStatus        OrderStatus = "NEW"
	InvalidStatus    OrderStatus = "INVALID"
	ProcessingStatus OrderStatus = "PROCESSING"
	ProcessedStatus  OrderStatus = "PROCESSED"
)

type Order struct {
	ID        int64
	OrderID   string
	UserID    int64
	Status    OrderStatus
	Accrual   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}
