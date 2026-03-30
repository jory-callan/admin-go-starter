package repo

import (
	"context"

	"aicode/internal/model"
	"aicode/pkg/goutils/gormutil"

	"gorm.io/gorm"
)

// UserRepo 用户仓库
type UserRepo struct {
	*gormutil.BaseRepo[model.User]
}

// NewUserRepo 创建用户仓库
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		BaseRepo: gormutil.NewBaseRepo[model.User](db),
	}
}

// GetByUsername 根据用户名查询
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.GetDB(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱查询
func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.GetDB(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByPhone 根据手机号查询
func (r *UserRepo) GetByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User
	err := r.GetDB(ctx).Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ListByRoleID 根据角色ID查询用户列表
func (r *UserRepo) ListByRoleID(ctx context.Context, roleID string) ([]model.User, error) {
	var users []model.User
	err := r.GetDB(ctx).
		Joins("JOIN user_roles ON user_roles.user_id = users.id").
		Where("user_roles.role_id = ?", roleID).
		Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
