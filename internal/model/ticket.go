package model

import (
	"time"

	"gorm.io/gorm"
)

// TicketStatus 工单状态
type TicketStatus string

const (
	TicketStatusPending  TicketStatus = "PENDING"
	TicketStatusLocked   TicketStatus = "LOCKED"
	TicketStatusExecuted TicketStatus = "EXECUTED"
	TicketStatusFailed   TicketStatus = "FAILED"
	TicketStatusRejected TicketStatus = "REJECTED"
)

// Ticket 工单模型
type Ticket struct {
	ID        string         `gorm:"primaryKey;type:varchar(36);comment:主键ID" json:"id"`
	CreatedAt time.Time      `gorm:"comment:创建时间" json:"created_at"`
	CreatedBy string         `gorm:"type:varchar(36);comment:创建人ID" json:"created_by"`
	UpdatedAt time.Time      `gorm:"comment:更新时间" json:"updated_at"`
	UpdatedBy string         `gorm:"type:varchar(36);comment:更新人ID" json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at"`
	DeletedBy string         `gorm:"type:varchar(36);comment:删除人ID" json:"deleted_by"`

	Title      string       `gorm:"type:varchar(200);not null;comment:工单标题" json:"title"`
	SQLContent string       `gorm:"type:text;not null;comment:SQL内容" json:"sql_content"`
	InstanceID string       `gorm:"not null;comment:关联实例ID" json:"instance_id"`
	DBName     string       `gorm:"type:varchar(100);not null;comment:目标数据库" json:"db_name"`
	Status     TicketStatus `gorm:"type:varchar(20);not null;default:'PENDING';comment:状态" json:"status"`
	CreatorID  string       `gorm:"not null;comment:创建人ID" json:"creator_id"`
	ExecutorID *string      `gorm:"null;comment:执行人ID" json:"executor_id"`
	ResultMsg  *string      `gorm:"type:text;null;comment:执行结果" json:"result_msg"`
	ExecutedAt *time.Time   `gorm:"null;comment:执行时间" json:"executed_at"`
}

func (Ticket) TableName() string {
	return "tickets"
}
