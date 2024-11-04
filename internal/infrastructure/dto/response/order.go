package response

import "gopher_mart/internal/domain"

type Order struct {
	OrderID   string             `json:"number"`
	Status    domain.OrderStatus `json:"status"`
	Accrual   int64              `json:"accrual"`
	CreatedAt string             `json:"uploaded_at"`
}
