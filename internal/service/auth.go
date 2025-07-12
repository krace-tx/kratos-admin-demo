package service

import (
	"context"

	pb "kratos-admin-demo/api/admin/v1"
)

type AuthService struct {
	pb.UnimplementedAuthServer
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	return &pb.RegisterReply{}, nil
}
func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	return &pb.LoginReply{}, nil
}
func (s *AuthService) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.LoginReply, error) {
	return &pb.LoginReply{}, nil
}
