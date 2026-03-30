package repo

import (
	"context"

	"aicode/internal/model"
	"aicode/pkg/goutils/gormutil"

	"gorm.io/gorm"
)

// InstanceRepo 数据库实例仓库
type InstanceRepo struct {
	*gormutil.BaseRepo[model.Instance]
}

// NewInstanceRepo 创建实例仓库
func NewInstanceRepo(db *gorm.DB) *InstanceRepo {
	return &InstanceRepo{
		BaseRepo: gormutil.NewBaseRepo[model.Instance](db),
	}
}

// GetByName 根据名称查询
func (r *InstanceRepo) GetByName(ctx context.Context, name string) (*model.Instance, error) {
	var instance model.Instance
	err := r.GetDB(ctx).Where("name = ?", name).First(&instance).Error
	if err != nil {
		return nil, err
	}
	return &instance, nil
}

// ListAll 获取所有实例
func (r *InstanceRepo) ListAll(ctx context.Context) ([]model.Instance, error) {
	var instances []model.Instance
	err := r.GetDB(ctx).Find(&instances).Error
	if err != nil {
		return nil, err
	}
	return instances, nil
}
