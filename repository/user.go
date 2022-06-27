package repository

import (
	"context"
	"rest-ws/models"
)

type UserRepository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	Close() error
}

var implementation UserRepository

func SetRepository(repository UserRepository) {
	implementation = repository
}

func InsertUser(ctx context.Context, user *models.User) error {
	return implementation.InsertUser(ctx, user)
}

func GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return implementation.GetUserByID(ctx, id)
}

func Close() error {
	return implementation.Close()
}