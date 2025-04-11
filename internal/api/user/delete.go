package user

import (
	"context"

	desc "github.com/Str1m/auth/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	if err := i.userService.Delete(ctx, req.GetId()); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete user")
	}
	return &emptypb.Empty{}, nil
}
