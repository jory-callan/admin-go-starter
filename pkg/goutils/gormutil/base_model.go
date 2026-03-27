package gormutil

// import (
// 	"time"

// 	"gorm.io/gorm"
// )

// // BaseModel 基础模型，提供一些通用字段
// type BaseModel struct {
// 	ID        string         `gorm:"type:char(36);primaryKey;comment:主键ID (UUID v7)" json:"id"`
// 	CreatedAt time.Time      `gorm:"index;autoCreateTime;comment:创建时间" json:"created_at"`
// 	CreatedBy string         `gorm:"type:char(36);comment:创建人ID" json:"created_by"`
// 	UpdatedAt time.Time      `gorm:"index;autoUpdateTime;comment:更新时间" json:"updated_at"`
// 	UpdatedBy string         `gorm:"type:char(36);comment:更新人ID" json:"updated_by"`
// 	DeletedAt gorm.DeletedAt `gorm:"index;comment:删除时间" json:"deleted_at"`
// 	DeletedBy string         `gorm:"type:char(36);comment:删除人ID" json:"deleted_by"`
// }
