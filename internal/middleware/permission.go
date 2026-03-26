package middleware

import (
	"aicode/internal/model"
	"aicode/pkg/jwt"
	"strings"
	"github.com/labstack/echo/v4"
)

// PermissionAuth 权限验证中间件
func PermissionAuth(requiredPermission string, permissions []string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 从上下文获取用户权限列表
			userPerms := c.Get("permissions").([]string)
			userID := c.Get("user_id").(string)
			
			// 检查是否有权限
			if !HasPermission(userPerms, requiredPermission) {
				return c.JSON(403, model.Response{
					Code: 403,
					Msg:  "无权限访问",
					Data: nil,
				})
			}
			
			// 传递用户ID给后续处理器
			c.Set("current_user_id", userID)
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
func HasPermission(userPerms []string, requiredPerm string) bool {
	for _, perm := range userPerms {
		// 用户有超级管理员权限
		if perm == "*" {
			return true
		}
		
		// 权限匹配（支持通配符）
		if matchPermission(perm, requiredPerm) {
			return true
		}
	}
	return false
}

// matchPermission 权限匹配（支持 * 通配符）
func matchPermission(userPerm, requiredPerm string) bool {
	// 完全匹配
	if userPerm == requiredPerm {
		return true
	}
	
	// 解析权限码
	userParts := strings.Split(userPerm, ":")
	reqParts := strings.Split(requiredPerm, ":")
	
	// 如果用户权限段数不足，无法匹配
	if len(userParts) < len(reqParts) {
		return false
	}
	
	// 逐段匹配，* 匹配任意值
	for i := 0; i < len(reqParts); i++ {
		userPart := userParts[i]
		reqPart := reqParts[i]
		
		if userPart == "*" {
			// 如果是 * 且是最后一段，或者后面都是 *，则匹配
			if i == len(userParts)-1 {
				return true
			}
			// 继续检查下一段
			continue
		}
		
		if userPart != reqPart {
			return false
		}
	}
	
	// 如果用户权限还有更多段，检查是否都是 *
	for i := len(reqParts); i < len(userParts); i++ {
		if userParts[i] != "*" {
			return false
		}
	}
	
	return true
}

// JWTAuth JWT认证中间件
func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// 从请求头获取 token
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(401, model.Response{
					Code: 401,
					Msg:  "未提供认证令牌",
					Data: nil,
				})
			}
			
			// 提取 Bearer token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(401, model.Response{
					Code: 401,
					Msg:  "认证令牌格式错误",
					Data: nil,
				})
			}
			
			token := parts[1]
			
			// 解析 token
			claims, err := jwt.ParseToken(token)
			if err != nil {
				return c.JSON(401, model.Response{
					Code: 401,
					Msg:  "认证令牌无效或已过期",
					Data: nil,
				})
			}
			
			// 将用户信息存入上下文
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("roles", claims.Roles)
			c.Set("permissions", claims.Permissions)
			
			return next(c)
		}
	}
}

// RequirePermission 权限验证装饰器（简化版本）
func RequirePermission(requiredPerm string) echo.MiddlewareFunc {
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
			
			if !HasPermission(userPerms, requiredPerm) {
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
