package model

type UserRole struct {
	UserID string `gorm:"primaryKey;type:varchar(36)" json:"user_id"`
	RoleID string `gorm:"primaryKey;type:varchar(36)" json:"role_id"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
