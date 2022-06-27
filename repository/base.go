package repository

import (
	"context"
	"rest-ws/models"
)

type ModelRepository[T models.Models] interface {
	Insert(ctx context.Context, data *T) (error, bool)
	GetById(ctx context.Context, id string) (*T, error)
	Close() error
}
