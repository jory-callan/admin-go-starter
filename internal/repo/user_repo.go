package repo

import (
	"context"
	"aicode/internal/model"
	"gorm.io/gorm"
)

// UserRepo 用户数据访问层
type UserRepo struct {
	*BaseRepo[model.User]
}

// NewUserRepo 创建用户 Repo
func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{
		BaseRepo: NewBaseRepo[model.User](db),
	}
}

// GetByUsername 根据用户名查询
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.DB.WithContext(ctx).
		Where("username = ?", username).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByIDWithRoles 根据 ID 查询用户（包含角色和权限）
func (r *UserRepo) GetByIDWithRoles(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.DB.WithContext(ctx).
		Where("id = ?", id).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// ListWithPagination 分页查询用户列表
func (r *UserRepo) ListWithPagination(ctx context.Context, query *model.PageQuery) (*model.PageResult[model.User], error) {
	var items []model.User
	var total int64

	db := r.DB.WithContext(ctx).Model(&model.User{})

	// 关键词搜索
	if query.Keyword != "" {
		db = db.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?",
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
		db = db.Order("created_at DESC")
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

	return &model.PageResult[model.User]{
		Items:   items,
		Total:   total,
		Page:    query.Page,
		Size:    query.Size,
		HasMore: hasMore,
	}, nil
}
