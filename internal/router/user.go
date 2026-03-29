package router

import (
	"aicode/internal/handler"

	"github.com/labstack/echo/v4"
)

// registerUserRoutes 注册用户管理路由
func (r *Router) registerUserRoutes(g *echo.Group) {
	h := handler.NewUserHandler(r.core.DB)

	users := g.Group("/users")
	{
		users.GET("", h.List)
		users.GET("/:id", h.Get)
		users.POST("", h.Create, r.mw.RequirePermission("system:user:create"))
		users.PUT("/:id", h.Update, r.mw.RequirePermission("system:user:update"))
		users.DELETE("/:id", h.Delete, r.mw.RequirePermission("system:user:delete"))
		users.PUT("/:id/roles", h.AssignRoles, r.mw.RequirePermission("system:user:assign"))
	}
}
