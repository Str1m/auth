package converter

import (
	modelService "github.com/Str1m/auth/internal/model"
	modelRepo "github.com/Str1m/auth/internal/storage/users/model"
)

func ToUserFromStorage(user *modelRepo.User) *modelService.User {
	return &modelService.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
