package router

import (
	"aicode/internal/handler"

	"github.com/labstack/echo/v4"
)

// registerRoleRoutes 注册角色管理路由
func (r *Router) registerRoleRoutes(g *echo.Group) {
	h := handler.NewRoleHandler(r.core.DB)

	roles := g.Group("/roles")
	{
		roles.GET("", h.List)
		roles.GET("/:id", h.Get)
		roles.POST("", h.Create, r.mw.RequirePermission("system:role:create"))
		roles.PUT("/:id", h.Update, r.mw.RequirePermission("system:role:update"))
		roles.DELETE("/:id", h.Delete, r.mw.RequirePermission("system:role:delete"))
		roles.PUT("/:id/permissions", h.AssignPermissions, r.mw.RequirePermission("system:role:assign"))
	}
}
