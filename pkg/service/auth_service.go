package service

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"

	"advanced-flight-server/pkg/database"
	"advanced-flight-server/pkg/entity"

	"gorm.io/gorm"
)

// 认证相关错误
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserInactive      = errors.New("user is inactive")
	ErrUserBanned        = errors.New("user is banned")
	ErrInvalidCredential = errors.New("invalid credential")
)

// AuthResult 认证结果
type AuthResult struct {
	Success bool
	User    *entity.User
	Error   error
}

// AuthService 认证服务接口
type AuthService interface {
	// ValidateLogin 验证用户登录
	// cid: 用户CID
	// password: 用户密码（明文）
	// 返回认证结果
	ValidateLogin(ctx context.Context, cid, password string) *AuthResult

	// ValidateLoginWithHash 验证用户登录（密码已哈希）
	// cid: 用户CID
	// passwordHash: 密码哈希值
	// 返回认证结果
	ValidateLoginWithHash(ctx context.Context, cid, passwordHash string) *AuthResult

	// GetUserByCID 根据CID获取用户信息
	GetUserByCID(ctx context.Context, cid string) (*entity.User, error)

	// CheckUserExists 检查用户是否存在
	CheckUserExists(ctx context.Context, cid string) (bool, error)
}

// authService 认证服务实现
type authService struct{}

// NewAuthService 创建认证服务
func NewAuthService() AuthService {
	return &authService{}
}

// db 获取账户数据库实例
func (s *authService) db() *gorm.DB {
	return database.Account()
}

// ValidateLogin 验证用户登录
func (s *authService) ValidateLogin(ctx context.Context, cid, password string) *AuthResult {
	// 查找用户
	user, err := s.findUserByCID(ctx, cid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &AuthResult{Success: false, Error: ErrUserNotFound}
		}
		return &AuthResult{Success: false, Error: err}
	}

	// 检查用户状态
	if result := s.checkUserStatus(user); result != nil {
		return result
	}

	// 验证密码
	if !s.verifyPassword(password, user.Password) {
		return &AuthResult{Success: false, User: nil, Error: ErrInvalidPassword}
	}

	return &AuthResult{Success: true, User: user, Error: nil}
}

// ValidateLoginWithHash 验证用户登录（密码已哈希）
func (s *authService) ValidateLoginWithHash(ctx context.Context, cid, passwordHash string) *AuthResult {
	// 查找用户
	user, err := s.findUserByCID(ctx, cid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &AuthResult{Success: false, Error: ErrUserNotFound}
		}
		return &AuthResult{Success: false, Error: err}
	}

	// 检查用户状态
	if result := s.checkUserStatus(user); result != nil {
		return result
	}

	// 直接比较哈希值
	if !s.compareHash(passwordHash, user.Password) {
		return &AuthResult{Success: false, User: nil, Error: ErrInvalidPassword}
	}

	return &AuthResult{Success: true, User: user, Error: nil}
}

// GetUserByCID 根据CID获取用户信息
func (s *authService) GetUserByCID(ctx context.Context, cid string) (*entity.User, error) {
	user, err := s.findUserByCID(ctx, cid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return user, nil
}

// CheckUserExists 检查用户是否存在
func (s *authService) CheckUserExists(ctx context.Context, cid string) (bool, error) {
	var count int64
	err := s.db().WithContext(ctx).Model(&entity.User{}).Where("cid = ?", cid).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// findUserByCID 根据CID查找用户
func (s *authService) findUserByCID(ctx context.Context, cid string) (*entity.User, error) {
	var user entity.User
	err := s.db().WithContext(ctx).Where("cid = ?", cid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// checkUserStatus 检查用户状态
func (s *authService) checkUserStatus(user *entity.User) *AuthResult {
	if user.IsBanned() {
		return &AuthResult{Success: false, User: user, Error: ErrUserBanned}
	}
	if !user.IsActive() {
		return &AuthResult{Success: false, User: user, Error: ErrUserInactive}
	}
	return nil
}

// verifyPassword 验证密码
// 这里使用SHA256作为示例，实际生产环境建议使用bcrypt或argon2
func (s *authService) verifyPassword(plainPassword, hashedPassword string) bool {
	hash := sha256.Sum256([]byte(plainPassword))
	computedHash := hex.EncodeToString(hash[:])
	return s.compareHash(computedHash, hashedPassword)
}

// compareHash 安全比较两个哈希值（防止时序攻击）
func (s *authService) compareHash(hash1, hash2 string) bool {
	return subtle.ConstantTimeCompare([]byte(hash1), []byte(hash2)) == 1
}
