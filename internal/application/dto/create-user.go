package dto

type CreateUser struct {
	Login        string `json:"login"`
	PasswordHash string `json:"passwordHash"`
}
