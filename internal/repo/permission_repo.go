package repo

import (
	"context"
	"aicode/internal/model"
	"gorm.io/gorm"
)

// PermissionRepo 权限数据访问层
type PermissionRepo struct {
	*BaseRepo[model.Permission]
}

// NewPermissionRepo 创建权限 Repo
func NewPermissionRepo(db *gorm.DB) *PermissionRepo {
	return &PermissionRepo{
		BaseRepo: NewBaseRepo[model.Permission](db),
	}
}

// GetByCode 根据权限码查询
func (r *PermissionRepo) GetByCode(ctx context.Context, code string) (*model.Permission, error) {
	var perm model.Permission
	err := r.DB.WithContext(ctx).Where("code = ?", code).First(&perm).Error
	if err != nil {
		return nil, err
	}
	return &perm, nil
}

// GetTree 获取权限树
func (r *PermissionRepo) GetTree(ctx context.Context) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.DB.WithContext(ctx).
		Where("status = ?", 1).
		Order("sort ASC, created_at ASC").
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return r.buildTree(permissions), nil
}

// buildTree 构建权限树
func (r *PermissionRepo) buildTree(permissions []model.Permission) []model.Permission {
	treeMap := make(map[string]*model.Permission)
	root := make([]model.Permission, 0)

	// 构建映射
	for i := range permissions {
		treeMap[permissions[i].ID] = &permissions[i]
	}

	// 构建树
	for i := range permissions {
		if permissions[i].ParentID == "" {
			root = append(root, permissions[i])
		} else {
			if parent, ok := treeMap[permissions[i].ParentID]; ok {
				parent.Children = append(parent.Children, permissions[i])
			}
		}
	}

	return root
}

// ListByRoleIDs 根据角色 ID 列表查询权限
func (r *PermissionRepo) ListByRoleIDs(ctx context.Context, roleIDs []string) ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.DB.WithContext(ctx).
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id IN ? AND permissions.status = ?", roleIDs, 1).
		Find(&permissions).Error
	return permissions, err
}

// ListWithPagination 分页查询权限列表
func (r *PermissionRepo) ListWithPagination(ctx context.Context, query *model.PageQuery) (*model.PageResult[model.Permission], error) {
	var items []model.Permission
	var total int64

	db := r.DB.WithContext(ctx).Model(&model.Permission{})

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
	if err := db.Offset(offset).Limit(query.Size).Find(&items).Error; err != nil {
		return nil, err
	}

	// 计算是否有下一页
	hasMore := false
	if query.NeedCount {
		hasMore = int64(query.Page*query.Size) < total
	}

	return &model.PageResult[model.Permission]{
		Items:   items,
		Total:   total,
		Page:    query.Page,
		Size:    query.Size,
		HasMore: hasMore,
	}, nil
}
