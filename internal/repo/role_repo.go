package repo

import (
	"context"

	"aicode/internal/model"
	"aicode/pkg/goutils/gormutil"

	"gorm.io/gorm"
)

// RoleRepo 角色仓库
type RoleRepo struct {
	*gormutil.BaseRepo[model.Role]
}

// NewRoleRepo 创建角色仓库
func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{
		BaseRepo: gormutil.NewBaseRepo[model.Role](db),
	}
}

// GetByCode 根据角色编码查询
func (r *RoleRepo) GetByCode(ctx context.Context, code string) (*model.Role, error) {
	var role model.Role
	err := r.GetDB(ctx).Where("code = ?", code).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// ListByUserID 根据用户ID查询角色列表
func (r *RoleRepo) ListByUserID(ctx context.Context, userID string) ([]model.Role, error) {
	var roles []model.Role
	err := r.GetDB(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	if err != nil {
		return nil, err
	}
	return roles, nil
}

// GetRoleCodesByUserID 获取用户的角色编码列表
func (r *RoleRepo) GetRoleCodesByUserID(ctx context.Context, userID string) ([]string, error) {
	var codes []string
	err := r.GetDB(ctx).
		Model(&model.Role{}).
		Select("roles.code").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Pluck("code", &codes).Error
	if err != nil {
		return nil, err
	}
	return codes, nil
}
