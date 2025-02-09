package repository

import (
	"context"
)

type BaseRepositoryInterface[T any] interface {
	Create(ctx context.Context, user *T) error
	FindByID(ctx context.Context, id uint) (*T, error)
	Update(ctx context.Context, user *T) error
	Delete(ctx context.Context, id uint) error
}
