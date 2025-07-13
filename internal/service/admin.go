package service

import (
	"admin/internal/biz"
	"context"

	pb "admin/api/admin/v1"
)

type AdminService struct {
	pb.UnimplementedAdminServer
	uc *biz.AdminUsecase
}

func NewAdminService(uc *biz.AdminUsecase) *AdminService {
	return &AdminService{uc: uc}
}

func (s *AdminService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	return &pb.RegisterReply{}, nil
}
func (s *AdminService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	return &pb.LoginReply{}, nil
}
func (s *AdminService) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.LoginReply, error) {
	return &pb.LoginReply{}, nil
}
