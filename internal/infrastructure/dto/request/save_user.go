package request

type SaveUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
	Name     string `json:"name"`
}