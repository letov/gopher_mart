package dto

type SaveUser struct {
	Login        string `json:"login"`
	PasswordHash string `json:"passwordHash"`
}
