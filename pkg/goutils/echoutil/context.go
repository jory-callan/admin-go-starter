package echoutil

import "github.com/labstack/echo/v4"

const (
	UserIDKey = "userID"
)

// SetUserID 设置当前请求中的用户 ID
func SetUserID(c echo.Context, userID string) {
	c.Set(UserIDKey, userID)
}

// GetUserID 获取当前请求中的用户 ID
func GetUserID(c echo.Context) string {
	if v := c.Get("user_id"); v != nil {
		return v.(string)
	}
	return ""
}

func GetUsername(c echo.Context) string {
	if v := c.Get("username"); v != nil {
		return v.(string)
	}
	return ""
}

func GetRoles(c echo.Context) []string {
	if v := c.Get("roles"); v != nil {
		return v.([]string)
	}
	return nil
}
