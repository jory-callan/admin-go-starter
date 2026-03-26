package service

import (
	"aicode/internal/model"
	"aicode/internal/repo"
	"context"
)

// UserService 用户服务
type UserService struct {
	userRepo *repo.UserRepo
	roleRepo *repo.RoleRepo
}

// NewUserService 创建用户服务
func NewUserService(userRepo *repo.UserRepo, roleRepo *repo.RoleRepo) *UserService {
	return &UserService{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

// List 用户列表
func (s *UserService) List(ctx context.Context, query *model.PageQuery) (*model.PageResult[model.User], error) {
	return s.userRepo.ListWithPagination(ctx, query)
}

// GetByID 获取用户详情
func (s *UserService) GetByID(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.GetByIDWithRoles(ctx, id)
}

// Create 创建用户
func (s *UserService) Create(ctx context.Context, req *CreateUserRequest, operatorID string) error {
	// 密码加密
	hashedPwd, err := HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &model.User{
		Username:  req.Username,
		Password:  hashedPwd,
		Nickname:  req.Nickname,
		Email:     req.Email,
		Phone:     req.Phone,
		Status:    1,
		CreatedBy: operatorID,
		UpdatedBy: operatorID,
	}

	// 分配角色
	if len(req.RoleIDs) > 0 {
		var roles []model.Role
		if err := s.userRepo.DB.WithContext(ctx).Where("id IN ?", req.RoleIDs).Find(&roles).Error; err != nil {
			return err
		}
		user.Roles = roles
	}

	return s.userRepo.Create(ctx, user)
}

// Update 更新用户
func (s *UserService) Update(ctx context.Context, id string, req *UpdateUserRequest, operatorID string) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}

	// 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Status != 0 {
		user.Status = req.Status
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	user.UpdatedBy = operatorID

	// 更新密码
	if req.Password != "" {
		hashedPwd, err := HashPassword(req.Password)
		if err != nil {
			return err
		}
		user.Password = hashedPwd
	}

	// 更新角色
	if len(req.RoleIDs) > 0 {
		var roles []model.Role
		if err := s.userRepo.DB.WithContext(ctx).Where("id IN ?", req.RoleIDs).Find(&roles).Error; err != nil {
			return err
		}
		user.Roles = roles
	}

	return s.userRepo.Update(ctx, user)
}

// Delete 删除用户
func (s *UserService) Delete(ctx context.Context, id, operatorID string) error {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	user.DeletedBy = operatorID
	return s.userRepo.Delete(ctx, id)
}

// AssignRoles 为用户分配角色
func (s *UserService) AssignRoles(ctx context.Context, userID string, roleIDs []string) error {
	var user model.User
	if err := s.userRepo.DB.WithContext(ctx).First(&user, "id = ?", userID).Error; err != nil {
		return err
	}

	return s.userRepo.DB.WithContext(ctx).
		Model(&user).
		Omit("Roles.*").
		Association("Roles").
		Replace(roleIDs)
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string   `json:"username" binding:"required"`
	Password string   `json:"password" binding:"required"`
	Nickname string   `json:"nickname"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	RoleIDs  []string `json:"role_ids"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	Nickname string   `json:"nickname"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Password string   `json:"password"`
	Avatar   string   `json:"avatar"`
	Status   int      `json:"status"`
	RoleIDs  []string `json:"role_ids"`
}
