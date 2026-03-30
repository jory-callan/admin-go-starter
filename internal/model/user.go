package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        string         `gorm:"primaryKey;type:varchar(36);comment:主键ID" json:"id"`
	CreatedAt time.Time      `gorm:"comment:创建时间" json:"created_at"`
	CreatedBy string         `gorm:"type:varchar(36);comment:创建人ID" json:"created_by"`
	UpdatedAt time.Time      `gorm:"comment:更新时间" json:"updated_at"`
	UpdatedBy string         `gorm:"type:varchar(36);comment:更新人ID" json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at"`
	DeletedBy string         `gorm:"type:varchar(36);comment:删除人ID" json:"deleted_by"`

	Username string `gorm:"type:varchar(50);not null;uniqueIndex" json:"username"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Email    string `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Phone    string `gorm:"type:varchar(20);uniqueIndex" json:"phone"`
	Status   int    `gorm:"type:tinyint;default:1;comment:状态:0禁用,1启用" json:"status"`
}
