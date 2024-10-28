package domain

import "time"

type Action string

func (a *Action) String() string {
	return string(*a)
}

const (
	SaveUserAction    Action = "SAVE_USER"
	LoginAction       Action = "LOGIN"
	CalcAccrualAction Action = "CALC_ACCRUAL"
)

type Event struct {
	ID        int64
	RootID    string
	Action    Action
	Payload   interface{}
	CreatedAt time.Time
}
