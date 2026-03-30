package repo

import (
	"context"
	"time"

	"aicode/internal/model"
	"aicode/pkg/goutils/gormutil"
	"aicode/pkg/goutils/response"

	"gorm.io/gorm"
)

// TicketRepo 工单仓库
type TicketRepo struct {
	*gormutil.BaseRepo[model.Ticket]
}

// NewTicketRepo 创建工单仓库
func NewTicketRepo(db *gorm.DB) *TicketRepo {
	return &TicketRepo{
		BaseRepo: gormutil.NewBaseRepo[model.Ticket](db),
	}
}

// ListByCreator 查询指定用户的工单
func (r *TicketRepo) ListByCreator(ctx context.Context, creatorID int64, pq *response.PageQuery) (*response.PageResult, error) {
	db := r.GetDB(ctx).Where("creator_id = ?", creatorID)
	return r.Pagination(ctx, pq, db)
}

// ListAll 查询所有工单
func (r *TicketRepo) ListAll(ctx context.Context, pq *response.PageQuery) (*response.PageResult, error) {
	return r.Pagination(ctx, pq, nil)
}

// GetByID 获取工单详情
func (r *TicketRepo) GetByID(ctx context.Context, id int64) (*model.Ticket, error) {
	var ticket model.Ticket
	err := r.GetDB(ctx).Where("id = ?", id).First(&ticket).Error
	if err != nil {
		return nil, err
	}
	return &ticket, nil
}

// UpdateStatus 更新工单状态
func (r *TicketRepo) UpdateStatus(ctx context.Context, id int64, status model.TicketStatus, executorID *int64, resultMsg *string) error {
	updates := map[string]interface{}{
		"status": status,
	}
	if executorID != nil {
		updates["executor_id"] = *executorID
	}
	if resultMsg != nil {
		updates["result_msg"] = *resultMsg
	}
	if status == model.TicketStatusExecuted || status == model.TicketStatusFailed {
		now := time.Now()
		updates["executed_at"] = &now
	}
	return r.GetDB(ctx).Model(&model.Ticket{}).Where("id = ?", id).Updates(updates).Error
}

// Execute 执行工单（更新状态为 LOCKED 然后执行）
func (r *TicketRepo) Execute(ctx context.Context, id int64, executorID int64) error {
	return r.GetDB(ctx).Transaction(func(tx *gorm.DB) error {
		// 先锁定工单
		if err := tx.Model(&model.Ticket{}).Where("id = ? AND status = ?", id, model.TicketStatusPending).
			Updates(map[string]interface{}{
				"status":      model.TicketStatusLocked,
				"executor_id": executorID,
			}).Error; err != nil {
			return err
		}
		return nil
	})
}
