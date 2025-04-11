package user

import (
	"context"

	desc "github.com/Str1m/auth/pkg/auth_v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	var name, email *string
	if req.GetName() != nil {
		name = &req.Name.Value
	}
	if req.GetEmail() != nil {
		email = &req.Email.Value
	}
	if err := i.userService.Update(ctx, req.GetId(), name, email); err != nil {
		return nil, status.Error(codes.Internal, "failed to update user")
	}
	return &emptypb.Empty{}, nil
}
