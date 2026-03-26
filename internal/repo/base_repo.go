package repo

import (
	"context"
	"aicode/internal/model"
	"gorm.io/gorm"
)

// BaseRepo 泛型基础 Repository，提供通用 CRUD 操作
type BaseRepo[T any] struct {
	DB *gorm.DB
}

// NewBaseRepo 创建基础 Repo
func NewBaseRepo[T any](db *gorm.DB) *BaseRepo[T] {
	return &BaseRepo[T]{DB: db}
}

// Create 创建记录
func (r *BaseRepo[T]) Create(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Create(entity).Error
}

// Update 更新记录
func (r *BaseRepo[T]) Update(ctx context.Context, entity *T) error {
	return r.DB.WithContext(ctx).Save(entity).Error
}

// Delete 软删除记录
func (r *BaseRepo[T]) Delete(ctx context.Context, id string) error {
	var entity T
	return r.DB.WithContext(ctx).Delete(&entity, "id = ?", id).Error
}

// GetByID 根据 ID 查询
func (r *BaseRepo[T]) GetByID(ctx context.Context, id string) (*T, error) {
	var entity T
	err := r.DB.WithContext(ctx).First(&entity, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &entity, nil
}

// List 列表查询
func (r *BaseRepo[T]) List(ctx context.Context, query *model.PageQuery) (*model.PageResult[T], error) {
	var items []T
	var total int64

	db := r.DB.WithContext(ctx).Model(new(T))

	// 排序
	if query.Order != "" {
		db = db.Order(query.Order)
	}

	// 分页
	offset := (query.Page - 1) * query.Size
	if err := db.Offset(offset).Limit(query.Size).Find(&items).Error; err != nil {
		return nil, err
	}

	// 是否统计总数
	if query.NeedCount {
		if err := db.Count(&total).Error; err != nil {
			return nil, err
		}
	}

	// 计算是否有下一页
	hasMore := false
	if query.NeedCount {
		hasMore = int64(query.Page*query.Size) < total
	}

	return &model.PageResult[T]{
		Items:   items,
		Total:   total,
		Page:    query.Page,
		Size:    query.Size,
		HasMore: hasMore,
	}, nil
}

// Where 条件查询（类似 mybatis-plus）
func (r *BaseRepo[T]) Where(conditions ...interface{}) *gorm.DB {
	return r.DB.Where(conditions[0], conditions[1:]...)
}

// Count 统计数量
func (r *BaseRepo[T]) Count(ctx context.Context, conditions ...interface{}) (int64, error) {
	var count int64
	db := r.DB.WithContext(ctx).Model(new(T))
	if len(conditions) > 0 {
		db = db.Where(conditions[0], conditions[1:]...)
	}
	err := db.Count(&count).Error
	return count, err
}
