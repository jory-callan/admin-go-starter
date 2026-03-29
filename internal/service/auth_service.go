package service

import (
	"aicode/internal/model"
	"aicode/internal/repo"
	"aicode/pkg/jwt"
	"context"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// AuthService 认证服务
type AuthService struct {
	userRepo *repo.UserRepo
	roleRepo *repo.RoleRepo
}

// NewAuthService 创建认证服务
func NewAuthService(db *gorm.DB) *AuthService {
	return &AuthService{
		userRepo: repo.NewUserRepo(db),
		roleRepo: repo.NewRoleRepo(db),
	}
}

// collectUserRolesAndPermissions 收集用户的角色和权限
func (s *AuthService) collectUserRolesAndPermissions(ctx context.Context, userID string) ([]string, []string, error) {
	// 查询用户的角色
	roles, err := s.roleRepo.GetUserRoles(ctx, userID)
	if err != nil {
		return nil, nil, err
	}

	roleCodes := make([]string, 0, len(roles))
	permissions := make(map[string]bool)

	for _, role := range roles {
		if role.Status != 1 {
			continue
		}
		roleCodes = append(roleCodes, role.Code)

		// 查询角色的权限
		perms, err := s.roleRepo.GetRolePermissions(ctx, role.ID)
		if err != nil {
			return nil, nil, err
		}
		for _, perm := range perms {
			if perm.Status == 1 {
				permissions[perm.Code] = true
			}
		}
	}

	permList := make([]string, 0, len(permissions))
	for perm := range permissions {
		permList = append(permList, perm)
	}

	return roleCodes, permList, nil
}

// Login 用户登录
func (s *AuthService) Login(ctx context.Context, username, password string) (*model.LoginResponse, error) {
	// 查询用户
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, ErrUserDisabled
	}

	// 收集角色和权限
	roles, permList, err := s.collectUserRolesAndPermissions(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	// 生成 JWT token
	token, err := jwt.GenerateToken(user.ID, user.Username, roles, permList)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Token: token,
		UserInfo: model.UserInfo{
			ID:          user.ID,
			Username:    user.Username,
			Nickname:    user.Nickname,
			Avatar:      user.Avatar,
			Email:       user.Email,
			Phone:       user.Phone,
			Roles:       roles,
			Permissions: permList,
		},
	}, nil
}

// GetUserInfo 获取用户信息
func (s *AuthService) GetUserInfo(ctx context.Context, userID string) (*model.UserInfo, error) {
	user, err := s.userRepo.GetByIDWithRoles(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 收集角色和权限
	roles, permList, err := s.collectUserRolesAndPermissions(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &model.UserInfo{
		ID:          user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Avatar:      user.Avatar,
		Email:       user.Email,
		Phone:       user.Phone,
		Roles:       roles,
		Permissions: permList,
	}, nil
}

// HashPassword 密码加密
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

var (
	ErrUserDisabled = &ServiceError{Code: 400, Msg: "用户已被禁用"}
	ErrInvalidCreds = &ServiceError{Code: 400, Msg: "用户名或密码错误"}
)

// ServiceError 服务错误
type ServiceError struct {
	Code int
	Msg  string
}

func (e *ServiceError) Error() string {
	return e.Msg
}
