package repo

import (
	"context"

	"aicode/internal/model"
	"aicode/pkg/goutils/gormutil"
	"aicode/pkg/goutils/response"

	"gorm.io/gorm"
)

// QueryHistoryRepo 查询历史仓库
type QueryHistoryRepo struct {
	*gormutil.BaseRepo[model.QueryHistory]
}

// NewQueryHistoryRepo 创建查询历史仓库
func NewQueryHistoryRepo(db *gorm.DB) *QueryHistoryRepo {
	return &QueryHistoryRepo{
		BaseRepo: gormutil.NewBaseRepo[model.QueryHistory](db),
	}
}

// Create 创建查询历史
func (r *QueryHistoryRepo) Create(ctx context.Context, history *model.QueryHistory) error {
	return r.GetDB(ctx).Create(history).Error
}

// ListByUser 分页查询指定用户的查询历史
func (r *QueryHistoryRepo) ListByUser(ctx context.Context, userID int64, pq *response.PageQuery) (*response.PageResult, error) {
	db := r.GetDB(ctx).Where("user_id = ?", userID).Order("created_at DESC")
	return r.Pagination(ctx, pq, db)
}

// ListByInstance 查询指定实例的查询历史
func (r *QueryHistoryRepo) ListByInstance(ctx context.Context, instanceID int64, pq *response.PageQuery) (*response.PageResult, error) {
	db := r.GetDB(ctx).Where("instance_id = ?", instanceID).Order("created_at DESC")
	return r.Pagination(ctx, pq, db)
}
