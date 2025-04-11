package storage

import (
	"context"
	"errors"
	modelService "github.com/Str1m/auth/internal/model"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type Repository interface {
	Create(ctx context.Context, info *modelService.UserInfo, hashedPassword []byte) (int64, error)
	Get(ctx context.Context, id int64) (*modelService.User, error)
	Update(ctx context.Context, id int64, name, email *string) error
	Delete(ctx context.Context, id int64) error
}
