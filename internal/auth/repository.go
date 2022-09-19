package auth

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, user *User) (string, error)
	Update(ctx context.Context, id string, user *User) (bool, error)
	Delete(ctx context.Context, id string) (bool, error)
	GetAll(ctx context.Context) ([]User, error)
	GetById(ctx context.Context, id string) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	CheckEmail(ctx context.Context, email string) (bool, error)
}
