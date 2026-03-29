package middleware

import (
	"aicode/internal/model"
	"aicode/pkg/jwt"
	"strings"

	"github.com/labstack/echo/v4"
)

// JWTAuth JWT认证中间件
func (m *Manager) JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(401, model.Response{
					Code: 401,
					Msg:  "未提供认证令牌",
					Data: nil,
				})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(401, model.Response{
					Code: 401,
					Msg:  "认证令牌格式错误",
					Data: nil,
				})
			}

			token := parts[1]

			claims, err := jwt.ParseToken(token)
			if err != nil {
				return c.JSON(401, model.Response{
					Code: 401,
					Msg:  "认证令牌无效或已过期",
					Data: nil,
				})
			}

			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)

			return next(c)
		}
	}
}
