package user

import (
	"context"

	"github.com/Str1m/auth/internal/converter"
	desc "github.com/Str1m/auth/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := i.userService.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user")
	}
	return &desc.GetResponse{User: converter.ToUserFromService(user)}, nil
}
