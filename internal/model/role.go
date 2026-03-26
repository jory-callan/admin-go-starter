package model

import (
	"time"

	"gorm.io/gorm"
)

// Role 角色表
type Role struct {
	ID        string         `gorm:"primaryKey;type:varchar(36);comment:主键ID" json:"id"`
	CreatedAt time.Time      `gorm:"comment:创建时间" json:"created_at"`
	CreatedBy string         `gorm:"type:varchar(36);comment:创建人ID" json:"created_by"`
	UpdatedAt time.Time      `gorm:"comment:更新时间" json:"updated_at"`
	UpdatedBy string         `gorm:"type:varchar(36);comment:更新人ID" json:"updated_by"`
	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at"`
	DeletedBy string         `gorm:"type:varchar(36);comment:删除人ID" json:"deleted_by"`

	Name        string       `gorm:"type:varchar(50);uniqueIndex;comment:角色名称" json:"name"`
	Code        string       `gorm:"type:varchar(50);uniqueIndex;comment:角色编码" json:"code"`
	Description string       `gorm:"type:varchar(200);comment:角色描述" json:"description"`
	Sort        int          `gorm:"type:int;default:0;comment:排序" json:"sort"`
	Status      int          `gorm:"type:tinyint;default:1;comment:状态(1:正常 2:禁用)" json:"status"`
	Users       []User       `gorm:"many2many:user_roles;comment:角色用户" json:"users,omitempty"`
	Permissions []Permission `gorm:"many2many:role_permissions;comment:角色权限" json:"permissions,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}
