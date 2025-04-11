package converter

import (
	"github.com/Str1m/auth/internal/model"
	desc "github.com/Str1m/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserFromService(user *model.User) *desc.User {
	return &desc.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}

func ToUserInfoFromService(userInfo *model.UserInfo) *desc.UserInfo {
	return &desc.UserInfo{
		Name:            userInfo.Name,
		Email:           userInfo.Email,
		Password:        userInfo.Password,
		PasswordConfirm: userInfo.PasswordConfirm,
		Role:            userInfo.Role,
	}
}

func ToUserInfoFromDesc(userInfo *desc.UserInfo) *model.UserInfo {
	return &model.UserInfo{
		Name:            userInfo.GetName(),
		Email:           userInfo.GetEmail(),
		Password:        userInfo.GetPassword(),
		PasswordConfirm: userInfo.GetPasswordConfirm(),
		Role:            userInfo.GetRole(),
	}
}
