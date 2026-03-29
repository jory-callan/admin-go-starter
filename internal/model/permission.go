package model

import (
	"time"

	"gorm.io/gorm"
)

// Permission 权限表
type Permission struct {
	ID        string         `gorm:"primaryKey;type:varchar(36);comment:主键ID" json:"id"`
	CreatedAt time.Time      `gorm:"comment:创建时间" json:"created_at"`
	CreatedBy string         `gorm:"type:varchar(36);comment:创建人ID" json:"created_by"`
	UpdatedAt time.Time      `gorm:"comment:更新时间" json:"updated_at"`
	UpdatedBy string         `gorm:"type:varchar(36);comment:更新人ID" json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at"`
	DeletedBy string         `gorm:"type:varchar(36);comment:删除人ID" json:"deleted_by"`

	Name        string       `gorm:"type:varchar(50);uniqueIndex;comment:权限名称" json:"name"`
	Code        string       `gorm:"type:varchar(100);uniqueIndex;comment:权限码(如system:user:write)" json:"code"`
	Description string       `gorm:"type:varchar(200);comment:权限描述" json:"description"`
	Type        int          `gorm:"type:tinyint;default:1;comment:权限类型(1:菜单 2:按钮 3:接口)" json:"type"`
	Sort        int          `gorm:"type:int;default:0;comment:排序" json:"sort"`
	Status      int          `gorm:"type:tinyint;default:1;comment:状态(1:正常 2:禁用)" json:"status"`
	ParentID string         `gorm:"type:varchar(36);comment:父级ID" json:"parent_id"`
	Children []Permission  `gorm:"-" json:"children,omitempty"`
}

func (Permission) TableName() string {
	return "permissions"
}
