package models

import (
	desc "github.com/Str1m/auth/pkg/auth_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserInfo struct {
	ID        int64
	Name      string
	Email     string
	Role      desc.Role
	CreatedAt *timestamppb.Timestamp
	UpdatedAt *timestamppb.Timestamp
}
