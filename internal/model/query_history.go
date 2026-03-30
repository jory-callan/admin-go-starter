package model

import (
	"time"

	"gorm.io/gorm"
)

// QueryHistory 查询历史记录
type QueryHistory struct {
	ID        string         `gorm:"primaryKey;type:varchar(36);comment:主键ID" json:"id"`
	CreatedAt time.Time      `gorm:"comment:创建时间" json:"created_at"`
	CreatedBy string         `gorm:"type:varchar(36);comment:创建人ID" json:"created_by"`
	UpdatedAt time.Time      `gorm:"comment:更新时间" json:"updated_at"`
	UpdatedBy string         `gorm:"type:varchar(36);comment:更新人ID" json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at"`
	DeletedBy string         `gorm:"type:varchar(36);comment:删除人ID" json:"deleted_by"`

	UserID       string `gorm:"not null;comment:执行人ID" json:"user_id"`
	InstanceID   string `gorm:"not null;comment:关联实例ID" json:"instance_id"`
	DBName       string `gorm:"type:varchar(100);not null;comment:目标数据库" json:"db_name"`
	SQLContent   string `gorm:"type:text;not null;comment:执行的SQL" json:"sql_content"`
	Duration     int64  `gorm:"type:bigint;comment:执行耗时(毫秒)" json:"duration"`
	RowsAffected int64  `gorm:"type:bigint;default:0;comment:影响行数" json:"rows_affected"`
	ErrorMsg     string `gorm:"type:text;comment:错误信息" json:"error_msg,omitempty"`
}

func (QueryHistory) TableName() string {
	return "query_histories"
}
