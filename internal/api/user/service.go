package user

import (
	"github.com/Str1m/auth/internal/service"

	desc "github.com/Str1m/auth/pkg/auth_v1"
)

type Implementation struct {
	desc.UnimplementedAuthV1Server
	userService service.Service
}

func NewImplementation(userService service.Service) *Implementation {
	return &Implementation{
		userService: userService,
	}
}
