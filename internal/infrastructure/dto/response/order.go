package response

type Status string

const (
	RegisteredStatus Status = "REGISTERED"
	InvalidStatus    Status = "INVALID"
	ProcessingStatus Status = "PROCESSING"
	ProcessedStatus  Status = "PROCESSED"
)

type Order struct {
	OrderID string `json:"order"`
	Status  Status `json:"status"`
	Accrual int64  `json:"accrual"`
}
