package user

import (
	"context"

	modelService "github.com/Str1m/auth/internal/model"
	desc "github.com/Str1m/auth/pkg/auth_v1"
)

type Service interface {
	Create(ctx context.Context, userInfo *modelService.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*modelService.User, error)
	Update(ctx context.Context, id int64, name *string, email *string) error
	Delete(ctx context.Context, id int64) error
}
type Implementation struct {
	desc.UnimplementedAuthV1Server
	userService Service
}

func NewImplementation(userService Service) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
