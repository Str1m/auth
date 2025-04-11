package user

import (
	"github.com/Str1m/auth/internal/client/db"
	"github.com/Str1m/auth/internal/storage"
	"log/slog"
)

type Service struct {
	log            *slog.Logger
	UserRepository storage.Repository
	TxManager      db.TxManager
}

func NewService(log *slog.Logger, authRepo storage.Repository, txManager db.TxManager) *Service {
	return &Service{
		log:            log,
		UserRepository: authRepo,
	}
}
