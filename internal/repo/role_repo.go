package repo

import (
	"context"
	"aicode/internal/model"
	"gorm.io/gorm"
)

// RoleRepo 角色数据访问层
type RoleRepo struct {
	*BaseRepo[model.Role]
}

// NewRoleRepo 创建角色 Repo
func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{
		BaseRepo: NewBaseRepo[model.Role](db),
	}
}

// GetByIDWithPermissions 根据 ID 查询角色（包含权限）
func (r *RoleRepo) GetByIDWithPermissions(ctx context.Context, id string) (*model.Role, error) {
	var role model.Role
	err := r.DB.WithContext(ctx).
		Preload("Permissions").
		Where("id = ?", id).
		First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetByCode 根据角色编码查询
func (r *RoleRepo) GetByCode(ctx context.Context, code string) (*model.Role, error) {
	var role model.Role
	err := r.DB.WithContext(ctx).
		Preload("Permissions").
		Where("code = ?", code).
		First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// ListWithPermissions 分页查询角色列表（包含权限）
func (r *RoleRepo) ListWithPermissions(ctx context.Context, query *model.PageQuery) (*model.PageResult[model.Role], error) {
	var items []model.Role
	var total int64

	db := r.DB.WithContext(ctx).Model(&model.Role{})

	// 关键词搜索
	if query.Keyword != "" {
		db = db.Where("name LIKE ? OR code LIKE ? OR description LIKE ?",
			"%"+query.Keyword+"%", "%"+query.Keyword+"%", "%"+query.Keyword+"%")
	}

	// 统计总数
	if query.NeedCount {
		if err := db.Count(&total).Error; err != nil {
			return nil, err
		}
	}

	// 排序
	if query.Order != "" {
		db = db.Order(query.Order)
	} else {
		db = db.Order("sort ASC, created_at DESC")
	}

	// 分页
	offset := (query.Page - 1) * query.Size
	if err := db.Offset(offset).Limit(query.Size).Preload("Permissions").Find(&items).Error; err != nil {
		return nil, err
	}

	// 计算是否有下一页
	hasMore := false
	if query.NeedCount {
		hasMore = int64(query.Page*query.Size) < total
	}

	return &model.PageResult[model.Role]{
		Items:   items,
		Total:   total,
		Page:    query.Page,
		Size:    query.Size,
		HasMore: hasMore,
	}, nil
}

// AssignPermissions 为角色分配权限
func (r *RoleRepo) AssignPermissions(ctx context.Context, roleID string, permissionIDs []string) error {
	return r.DB.WithContext(ctx).
		Model(&model.Role{}).
		Where("id = ?", roleID).
		Omit("Permissions.*").
		Association("Permissions").
		Replace(permissionIDs)
}

// GetUserRoles 获取用户的角色
func (r *RoleRepo) GetUserRoles(ctx context.Context, userID string) ([]model.Role, error) {
	var roles []model.Role
	err := r.DB.WithContext(ctx).
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Where("user_roles.user_id = ?", userID).
		Find(&roles).Error
	return roles, err
}
