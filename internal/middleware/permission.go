package middleware

import (
	"aicode/internal/model"
	"strings"

	"github.com/labstack/echo/v4"
)

// PermissionAuth 权限验证中间件
func (m *Manager) PermissionAuth(requiredPermission string, permissions []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userPerms := c.Get("permissions").([]string)
			userID := c.Get("user_id").(string)

			if !m.HasPermission(userPerms, requiredPermission) {
				return c.JSON(403, model.Response{
					Code: 403,
					Msg:  "无权限访问",
					Data: nil,
				})
			}

			c.Set("current_user_id", userID)
			return next(c)
		}
	}
}

// RequirePermission 权限验证装饰器（简化版本）
func (m *Manager) RequirePermission(requiredPerm string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userPerms, ok := c.Get("permissions").([]string)
			if !ok {
				return c.JSON(403, model.Response{
					Code: 403,
					Msg:  "无法获取用户权限",
					Data: nil,
				})
			}

			if !m.HasPermission(userPerms, requiredPerm) {
				return c.JSON(403, model.Response{
					Code: 403,
					Msg:  "无权限访问",
					Data: nil,
				})
			}

			return next(c)
		}
	}
}

// HasPermission 检查用户是否有指定权限（支持通配符）
// 规则：
// 1. * 代表所有权限（超级管理员）
// 2. system:* 代表系统管理所有权限
// 3. system:user:* 代表用户模块所有权限
// 4. system:user:write 代表用户模块写入权限
func (m *Manager) HasPermission(userPerms []string, requiredPerm string) bool {
	for _, perm := range userPerms {
		if perm == "*" {
			return true
		}
		if matchPermission(perm, requiredPerm) {
			return true
		}
	}
	return false
}

// matchPermission 权限匹配（支持 * 通配符）
func matchPermission(userPerm, requiredPerm string) bool {
	if userPerm == requiredPerm {
		return true
	}

	userParts := strings.Split(userPerm, ":")
	reqParts := strings.Split(requiredPerm, ":")

	if len(userParts) < len(reqParts) {
		return false
	}

	for i := 0; i < len(reqParts); i++ {
		userPart := userParts[i]
		reqPart := reqParts[i]

		if userPart == "*" {
			if i == len(userParts)-1 {
				return true
			}
			continue
		}

		if userPart != reqPart {
			return false
		}
	}

	for _, part := range userParts[len(reqParts):] {
		if part != "*" {
			return false
		}
	}

	return true
}
