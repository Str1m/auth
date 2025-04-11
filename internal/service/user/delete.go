package user

import (
	"context"
	"fmt"
)

func (s *Service) Delete(ctx context.Context, id int64) error {
	const op = "service.user.Delete"
	err := s.UserDBClient.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
