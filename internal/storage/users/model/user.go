package model

import (
	"time"

	desc "github.com/Str1m/auth/pkg/auth_v1"
)

type User struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	Role      desc.Role `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserInfo struct {
	Name            string    `db:"id"`
	Email           string    `db:"email"`
	Password        string    `db:"password"`
	PasswordConfirm string    `db:"password_confirm"`
	Role            desc.Role `db:"role"`
}
