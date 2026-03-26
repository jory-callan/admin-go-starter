package model

// PageQuery 分页查询参数
type PageQuery struct {
	Page      int    `json:"page" query:"page"`             // 页码，从 1 开始
	Size      int    `json:"size" query:"size"`             // 每页数量
	NeedCount bool   `json:"need_count" query:"need_count"` // 是否需要总数
	Order     string `json:"order" query:"order"`           // 排序字段，如 "created_at desc"
	Keyword   string `json:"keyword" query:"keyword"`       // 关键词搜索
}

// PageResult 分页返回结果
type PageResult[T any] struct {
	Items   []T   `json:"items"`   // 数据列表
	Total   int64 `json:"total"`   // 总记录数
	Page    int   `json:"page"`    // 当前页码
	Size    int   `json:"size"`    // 每页数量
	HasMore bool  `json:"hasMore"` // 是否有下一页
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	Token    string      `json:"token"`
	UserInfo UserInfo    `json:"user_info"`
}

// UserInfo 用户信息（不含密码）
type UserInfo struct {
	ID       string   `json:"id"`
	Username string   `json:"username"`
	Nickname string   `json:"nickname"`
	Avatar   string   `json:"avatar"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Roles    []string `json:"roles"`     // 角色编码列表
	Permissions []string `json:"permissions"` // 权限码列表
}
