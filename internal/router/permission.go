package router

import (
	"aicode/internal/handler"

	"github.com/labstack/echo/v4"
)

// registerPermissionRoutes 注册权限管理路由
func (r *Router) registerPermissionRoutes(g *echo.Group) {
	h := handler.NewPermissionHandler(r.core.DB)

	permissions := g.Group("/permissions")
	{
		permissions.GET("", h.List)
		permissions.GET("/tree", h.GetTree)
		permissions.GET("/:id", h.Get)
		permissions.POST("", h.Create, r.mw.RequirePermission("system:permission:create"))
		permissions.PUT("/:id", h.Update, r.mw.RequirePermission("system:permission:update"))
		permissions.DELETE("/:id", h.Delete, r.mw.RequirePermission("system:permission:delete"))
	}
}
