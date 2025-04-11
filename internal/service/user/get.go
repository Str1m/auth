package user

import (
	"context"
	"fmt"

	modelService "github.com/Str1m/auth/internal/model"
)

func (s *Service) Get(ctx context.Context, id int64) (*modelService.User, error) {
	const op = "service.user.Get"
	user, err := s.UserDBClient.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}
