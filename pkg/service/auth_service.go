package service

import (
	"context"
	"errors"
	"strconv"

	"advanced-flight-server/pkg/database"
	"advanced-flight-server/pkg/entity"
	"advanced-flight-server/pkg/protocol/pdu"

	"gorm.io/gorm"
)

// 认证相关错误
var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserInactive      = errors.New("user is inactive")
	ErrUserBanned        = errors.New("user is banned")
	ErrInvalidCredential = errors.New("invalid credential")
	ErrRatingTooHigh     = errors.New("requested rating exceeds user level")
)

// AuthResult 认证结果
type AuthResult struct {
	Success bool
	User    *entity.User
	Level   pdu.NetworkRating // 用户的等级
	Error   error
}

// AuthService 认证服务接口
type AuthService interface {
	// ValidateLogin 验证用户登录
	// cid: 用户CID
	// password: 用户密码（明文）
	// 返回认证结果
	ValidateLogin(ctx context.Context, cid, password string) *AuthResult

	// ValidateLoginWithRating 验证用户登录并检查rating
	// cid: 用户CID
	// password: 用户密码（明文）
	// requestedRating: 请求使用的等级
	// 返回认证结果，登录时请求的rating必须 <= 用户的level
	ValidateLoginWithRating(ctx context.Context, cid, password string, requestedRating pdu.NetworkRating) *AuthResult

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
func (s *authService) db() (*gorm.DB, error) {
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

	// 验证密码（明文比较）
	if password != user.Password {
		return &AuthResult{Success: false, User: nil, Error: ErrInvalidPassword}
	}

	level := s.getUserLevel(user)
	return &AuthResult{Success: true, User: user, Level: level, Error: nil}
}

// ValidateLoginWithRating 验证用户登录并检查rating
func (s *authService) ValidateLoginWithRating(ctx context.Context, cid, password string, requestedRating pdu.NetworkRating) *AuthResult {
	// 先进行基本登录验证
	result := s.ValidateLogin(ctx, cid, password)
	if !result.Success {
		return result
	}

	// 检查请求的rating是否超过用户level
	if requestedRating > result.Level {
		return &AuthResult{Success: false, User: result.User, Level: result.Level, Error: ErrRatingTooHigh}
	}

	return result
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
	db, err := s.db()
	if err != nil {
		return false, err
	}
	var count int64
	err = db.WithContext(ctx).Model(&entity.User{}).Where("username = ?", cid).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// findUserByCID 根据CID查找用户
func (s *authService) findUserByCID(ctx context.Context, cid string) (*entity.User, error) {
	db, err := s.db()
	if err != nil {
		return nil, err
	}
	var user entity.User
	err = db.WithContext(ctx).Where("username = ?", cid).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// checkUserStatus 检查用户状态
func (s *authService) checkUserStatus(user *entity.User) *AuthResult {
	// level为0表示被ban
	if s.isBanned(user) {
		return &AuthResult{Success: false, User: user, Error: ErrUserBanned}
	}
	// IsPassExam为0表示不active
	if !s.isActive(user) {
		return &AuthResult{Success: false, User: user, Error: ErrUserInactive}
	}
	return nil
}

// isBanned 检查用户是否被禁用 (level为0表示被ban)
func (s *authService) isBanned(user *entity.User) bool {
	level, err := strconv.Atoi(user.Level)
	if err != nil {
		return false
	}
	return level == 0
}

// isActive 检查用户是否激活 (IsPassExam不为0表示active)
func (s *authService) isActive(user *entity.User) bool {
	if user.IsPassExam == nil {
		return false
	}
	return *user.IsPassExam != 0
}

// getUserLevel 获取用户等级
func (s *authService) getUserLevel(user *entity.User) pdu.NetworkRating {
	level, err := strconv.Atoi(user.Level)
	if err != nil {
		return pdu.RatingOBS
	}
	return pdu.NetworkRating(level)
}
