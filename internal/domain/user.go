package domain

import "context"

type User struct {
	Login        string
	PasswordHash string
}

type UserRepository interface {
	Save(ctx context.Context, u User) error
}
