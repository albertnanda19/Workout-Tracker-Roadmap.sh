package domain

import "context"

type User struct {
	ID    int64
	Email string
	Name  string
}

type UserRepository interface {
	Create(ctx context.Context, entity *User) error
}
