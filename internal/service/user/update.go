package user

import (
	"context"
	"fmt"
)

func (s *Service) Update(ctx context.Context, id int64, name *string, email *string) error {
	const op = "service.user.Update"
	err := s.UserDBClient.Update(ctx, id, name, email)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
