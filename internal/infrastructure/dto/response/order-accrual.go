package response

type Status string

const (
	RegisteredStatus Status = "REGISTERED"
	InvalidStatus    Status = "INVALID"
	ProcessingStatus Status = "PROCESSING"
	ProcessedStatus  Status = "PROCESSED"
)

type OrderAccrual struct {
	OrderID int64  `json:"order"`
	Status  Status `json:"status"`
	Accrual int64  `json:"accrual"`
}
