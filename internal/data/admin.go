package data

import (
	"admin/internal/biz"
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"time"

	"github.com/google/uuid"
)

type userRepo struct {
	data *Data // 数据源，如 db 或 redis 实例
	log  *log.Helper
}

// NewUserRepo 构造函数
func NewUserRepo(data *Data, logger log.Logger) biz.AdminRepo {
	return &userRepo{data: data, log: log.NewHelper(logger)}
}

// 模拟数据库表结构
type UserModel struct {
	ID       string
	Username string
	Password string
	Email    string
}

// Register 用户注册
func (r *userRepo) Register(ctx context.Context, u *biz.User, emailCode string) (string, error) {
	// TODO: 验证 emailCode 是否正确，可从 Redis 查询
	if emailCode != "123456" {
		return "", fmt.Errorf("邮箱验证码无效")
	}

	// 生成唯一 ID
	uid := uuid.NewString()

	// 保存到数据库（这里只是模拟）
	user := &UserModel{
		ID:       uid,
		Username: u.Username,
		Password: u.Password, // TODO: 应该 hash
		Email:    u.Email,
	}

	// 假设 data.db 是 gorm.DB
	if err := r.data.db.WithContext(ctx).Create(user).Error; err != nil {
		return "", err
	}

	return uid, nil
}

// Login 用户登录
func (r *userRepo) Login(ctx context.Context, username, password, vCode string) (*biz.AuthToken, error) {
	// TODO: 验证行为验证码（vCode）
	if vCode != "pass" {
		return nil, fmt.Errorf("行为验证码无效")
	}

	var user UserModel
	if err := r.data.db.WithContext(ctx).
		Where("username = ? AND password = ?", username, password). // TODO: 应该加 hash 对比
		First(&user).Error; err != nil {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	// 模拟生成 token
	return &biz.AuthToken{
		AccessToken:  "access-token-" + user.ID,
		RefreshToken: "refresh-token-" + user.ID,
		Exp:          time.Now().Add(2 * time.Hour).Unix(),
	}, nil
}

// RefreshToken 刷新 AccessToken
func (r *userRepo) RefreshToken(ctx context.Context, refreshToken string) (*biz.AuthToken, error) {
	// TODO: 验证 refreshToken 有效性
	if refreshToken == "" {
		return nil, fmt.Errorf("无效的刷新令牌")
	}

	// 模拟从 refreshToken 获取 userId
	userId := refreshToken[len("refresh-token-"):] // 不严谨，仅演示

	return &biz.AuthToken{
		AccessToken:  "access-token-" + userId,
		RefreshToken: refreshToken,
		Exp:          time.Now().Add(2 * time.Hour).Unix(),
	}, nil
}
