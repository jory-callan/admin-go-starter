package repo

import (
	"context"

	"aicode/internal/model"

	"gorm.io/gorm"
)

// UserRoleRepo 用户角色关系仓库
type UserRoleRepo struct {
	db *gorm.DB
}

// NewUserRoleRepo 创建用户角色关系仓库
func NewUserRoleRepo(db *gorm.DB) *UserRoleRepo {
	return &UserRoleRepo{db: db}
}

// AddRoleForUser 为用户添加角色
func (r *UserRoleRepo) AddRoleForUser(ctx context.Context, userID, roleID string) error {
	userRole := &model.UserRole{
		UserID: userID,
		RoleID: roleID,
	}
	return r.db.WithContext(ctx).Create(userRole).Error
}

// RemoveRoleFromUser 从用户移除角色
func (r *UserRoleRepo) RemoveRoleFromUser(ctx context.Context, userID, roleID string) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Delete(&model.UserRole{}).Error
}

// SetRolesForUser 设置用户的角色（先删后加）
func (r *UserRoleRepo) SetRolesForUser(ctx context.Context, userID string, roleIDs []string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 先删除用户的所有角色
		if err := tx.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
			return err
		}
		// 再添加新角色
		if len(roleIDs) == 0 {
			return nil
		}
		userRoles := make([]*model.UserRole, len(roleIDs))
		for i, roleID := range roleIDs {
			userRoles[i] = &model.UserRole{
				UserID: userID,
				RoleID: roleID,
			}
		}
		return tx.CreateInBatches(userRoles, 100).Error
	})
}

// GetUserRoles 获取用户的角色ID列表
func (r *UserRoleRepo) GetUserRoles(ctx context.Context, userID string) ([]string, error) {
	var roleIDs []string
	err := r.db.WithContext(ctx).
		Model(&model.UserRole{}).
		Where("user_id = ?", userID).
		Pluck("role_id", &roleIDs).Error
	if err != nil {
		return nil, err
	}
	return roleIDs, nil
}

// GetUserRoleMap 获取多个用户的角色映射关系
func (r *UserRoleRepo) GetUserRoleMap(ctx context.Context, userIDs []string) (map[string][]string, error) {
	var userRoles []model.UserRole
	err := r.db.WithContext(ctx).
		Where("user_id IN ?", userIDs).
		Find(&userRoles).Error
	if err != nil {
		return nil, err
	}

	result := make(map[string][]string)
	for _, ur := range userRoles {
		result[ur.UserID] = append(result[ur.UserID], ur.RoleID)
	}
	return result, nil
}

// ExistsUserRole 检查用户角色关系是否存在
func (r *UserRoleRepo) ExistsUserRole(ctx context.Context, userID, roleID string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&model.UserRole{}).
		Where("user_id = ? AND role_id = ?", userID, roleID).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
