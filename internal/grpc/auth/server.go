package auth

import (
	"context"

	desc "github.com/Str1m/auth/pkg/auth_v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Server struct {
	desc.UnimplementedAuthV1Server
	repo Repo
}

func New(repo Repo) *Server {
	return &Server{
		repo: repo,
	}
}

type Repo interface {
	Create(ctx context.Context, info *desc.UserInfo) (int64, error)
	Get(ctx context.Context, id int64) (*desc.User, error)
	Update(ctx context.Context, id int64, name, email *string) error
	Delete(ctx context.Context, id int64) error
}

func (s *Server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	id, err := s.repo.Create(ctx, req.GetUserInfo())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to create user")
	}
	return &desc.CreateResponse{
		Id: id,
	}, nil
}

func (s *Server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	user, err := s.repo.Get(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get user")
	}
	return &desc.GetResponse{User: user}, nil
}

func (s *Server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	var name, email *string
	if req.GetName() != nil {
		name = &req.Name.Value
	}
	if req.GetEmail() != nil {
		email = &req.Email.Value
	}
	if err := s.repo.Update(ctx, req.GetId(), name, email); err != nil {
		return nil, status.Error(codes.Internal, "failed to update user")
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	if err := s.repo.Delete(ctx, req.GetId()); err != nil {
		return nil, status.Error(codes.Internal, "failed to delete user")
	}
	return &emptypb.Empty{}, nil
}
