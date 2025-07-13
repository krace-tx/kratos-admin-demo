package biz

import (
	"context"
	"github.com/go-kratos/kratos/v2/log"
)

// User 实体结构
type User struct {
	Username string
	Password string
	Email    string
}

// AuthToken 登录返回的 Token 信息
type AuthToken struct {
	AccessToken  string
	RefreshToken string
	Exp          int64
}

// AdminRepo 定义用户仓储接口
type AdminRepo interface {
	Register(ctx context.Context, u *User, emailCode string) (string, error)
	Login(ctx context.Context, username, password, vCode string) (*AuthToken, error)
	RefreshToken(ctx context.Context, refreshToken string) (*AuthToken, error)
}

// AdminUsecase 业务用例结构体
type AdminUsecase struct {
	repo AdminRepo
	log  *log.Helper
}

// NewAdminUsecase 构造函数
func NewAdminUsecase(repo AdminRepo, logger log.Logger) *AdminUsecase {
	return &AdminUsecase{repo: repo, log: log.NewHelper(logger)}
}

// 注册用户
func (uc *AdminUsecase) Register(ctx context.Context, u *User, emailCode string) (string, error) {
	uc.log.WithContext(ctx).Infof("Register: %v", u.Username)
	return uc.repo.Register(ctx, u, emailCode)
}

// 登录用户
func (uc *AdminUsecase) Login(ctx context.Context, username, password, vCode string) (*AuthToken, error) {
	uc.log.WithContext(ctx).Infof("Login: %v", username)
	return uc.repo.Login(ctx, username, password, vCode)
}

// 刷新 Token
func (uc *AdminUsecase) Refresh(ctx context.Context, refreshToken string) (*AuthToken, error) {
	uc.log.WithContext(ctx).Infof("Refresh: %v", refreshToken)
	return uc.repo.RefreshToken(ctx, refreshToken)
}
