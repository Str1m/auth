package model

import (
	"time"

	desc "github.com/Str1m/auth/pkg/auth_v1"
)

type User struct {
	ID        int64
	Name      string
	Email     string
	Role      desc.Role
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserInfo struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            desc.Role
}
