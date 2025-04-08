package auth

import (
	"log/slog"

	"github.com/Str1m/auth/internal/models"
)

type Repo interface {
	SaveUser(email string, passHash []byte) (int64, error)
	UpdateUser(id int64, name, email *string) error
	User(email string) (models.UserInfo, error)
	DeleteUser(id int64) error
}

type Auth struct {
	log      *slog.Logger
	authRepo Repo
}

func New(log *slog.Logger, authRepo Repo) *Auth {
	return &Auth{
		log:      log,
		authRepo: authRepo,
	}
}

// func Create(ctx context.Context, name, email, password, passwordConfirm string, role desc.Role) (int64, error) {
// 	panic("implement me")
// }

// func Get(ctx context.Context, id int64) (models.UserInfo, error) {
// 	panic("implement me")
// }

// func Update(ctx context.Context, id int64, name, email *string) error {
// 	panic("implement me")
// }

// func Delete(ctx context.Context, id int64) error {
// 	panic("implement me")
// }
