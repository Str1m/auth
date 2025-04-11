package service

import (
	"errors"
)

var (
	ErrPassNotEqual = errors.New("passwords are not equal")
)

//type Service interface {
//	Create(ctx context.Context, userInfo *modelService.UserInfo) (int64, error)
//	Get(ctx context.Context, id int64) (*modelService.User, error)
//	Update(ctx context.Context, id int64, name *string, email *string) error
//	Delete(ctx context.Context, id int64) error
//}
