package user

import (
	"context"
	"log/slog"

	modelService "github.com/Str1m/auth/internal/model"
)

type Repository interface {
	Create(ctx context.Context, info *modelService.UserInfo, hashedPassword []byte) (int64, error)
	Get(ctx context.Context, id int64) (*modelService.User, error)
	Update(ctx context.Context, id int64, name, email *string) error
	Delete(ctx context.Context, id int64) error
}

type Service struct {
	log            *slog.Logger
	UserRepository Repository
}

func NewService(log *slog.Logger, authRepo Repository) *Service {
	return &Service{
		log:            log,
		UserRepository: authRepo,
	}
}
