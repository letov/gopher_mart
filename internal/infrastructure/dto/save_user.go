package dto

type SaveUser struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}
