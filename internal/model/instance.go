package model

import (
	"time"

	"gorm.io/gorm"
)

// Instance 数据库实例模型
type Instance struct {
	ID        string         `gorm:"primaryKey;type:varchar(36);comment:主键ID" json:"id"`
	CreatedAt time.Time      `gorm:"comment:创建时间" json:"created_at"`
	CreatedBy string         `gorm:"type:varchar(36);comment:创建人ID" json:"created_by"`
	UpdatedAt time.Time      `gorm:"comment:更新时间" json:"updated_at"`
	UpdatedBy string         `gorm:"type:varchar(36);comment:更新人ID" json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at"`
	DeletedBy string         `gorm:"type:varchar(36);comment:删除人ID" json:"deleted_by"`

	Name      string `gorm:"type:varchar(100);not null;comment:实例名称" json:"name"`
	Host      string `gorm:"type:varchar(100);not null;comment:主机地址" json:"host"`
	Port      int    `gorm:"type:int;not null;comment:端口" json:"port"`
	AdminUser string `gorm:"type:varchar(50);not null;comment:管理员用户名" json:"admin_user"`
	AdminPass string `gorm:"type:varchar(255);not null;comment:管理员密码" json:"admin_pass"`
}

func (Instance) TableName() string {
	return "instances"
}
