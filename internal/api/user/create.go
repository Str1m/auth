package user

import (
	"context"

	"github.com/Str1m/auth/internal/converter"
	desc "github.com/Str1m/auth/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := i.userService.Create(ctx, converter.ToUserInfoFromDesc(req.GetUserInfo()))
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}
	return &desc.CreateResponse{
		Id: id,
	}, nil
}
