package response

type Balance struct {
	Current   int64 `json:"current"`
	Withdrawn int64 `json:"withdrawn"`
}
