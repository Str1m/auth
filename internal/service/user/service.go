package user

import (
	"context"
	modelService "github.com/Str1m/auth/internal/model"
	"log/slog"
)

type DBLayer interface {
	Create(ctx context.Context, info *modelService.UserInfo, hashedPassword []byte) (int64, error)
	Get(ctx context.Context, id int64) (*modelService.User, error)
	Update(ctx context.Context, id int64, name, email *string) error
	Delete(ctx context.Context, id int64) error
}
type Service struct {
	log          *slog.Logger
	UserDBClient DBLayer
	//TxManager      db.TxManager
}

func NewService(log *slog.Logger, dbClient DBLayer) *Service {
	return &Service{
		log:          log,
		UserDBClient: dbClient,
	}
}
