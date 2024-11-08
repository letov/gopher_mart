package domain

import "time"

type Action string

func (a *Action) String() string {
	return string(*a)
}

const (
	SaveUserAction      Action = "SAVE_USER"
	LoginAction         Action = "LOGIN"
	SaveOrderAction     Action = "SAVE_ORDER"
	UpdateOrderAction   Action = "UPDATE_ORDER"
	SaveOperationAction Action = "SAVE_OPERATION"
)

type Event struct {
	ID        int64
	RootID    string
	Action    Action
	Payload   interface{}
	CreatedAt time.Time
}
