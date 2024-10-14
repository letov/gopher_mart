package domain

type Action string

const (
	SaveUserAction Action = "SaveUserAction"
	LoginAction    Action = "LoginAction"
)

type Event struct {
	ID      int64
	RootID  string
	Action  Action
	Payload interface{}
}
