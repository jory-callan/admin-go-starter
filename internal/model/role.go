package model

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID        string         `gorm:"primaryKey;type:varchar(36);comment:主键ID" json:"id"`
	CreatedAt time.Time      `gorm:"comment:创建时间" json:"created_at"`
	CreatedBy string         `gorm:"type:varchar(36);comment:创建人ID" json:"created_by"`
	UpdatedAt time.Time      `gorm:"comment:更新时间" json:"updated_at"`
	UpdatedBy string         `gorm:"type:varchar(36);comment:更新人ID" json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at"`
	DeletedBy string         `gorm:"type:varchar(36);comment:删除人ID" json:"deleted_by"`

	Code string `gorm:"type:varchar(50);not null;uniqueIndex" json:"code"`
	Name string `gorm:"type:varchar(50);not null" json:"name"`
}

func (Role) TableName() string {
	return "roles"
}

/*
-- 初始化数据：
-- 1: admin (全权：创建实例、执行工单、查看所有)
-- 2: auditor (审计：执行工单、查看所有)
-- 3: developer (开发：仅查询、提交工单、查看自己工单)
*/
