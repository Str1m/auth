package user

import (
	"context"
	"fmt"

	"github.com/Str1m/auth/internal/service"

	modelService "github.com/Str1m/auth/internal/model"
	"golang.org/x/crypto/bcrypt"
)

func (s *Service) Create(ctx context.Context, userInfo *modelService.UserInfo) (int64, error) {
	const op = "service.user.Create"
	if userInfo.Password != userInfo.PasswordConfirm {
		return 0, fmt.Errorf("%s: failed to hash password: %w", op, service.ErrPassNotEqual)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userInfo.Password), bcrypt.DefaultCost)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to hash password: %w", op, err)
	}

	userInfo.Password, userInfo.PasswordConfirm = "", ""

	id, err := s.UserRepository.Create(ctx, userInfo, hashedPassword)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}
