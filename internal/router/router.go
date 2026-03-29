package router

import (
	"aicode/internal/app/core"
	"aicode/internal/middleware"

	"github.com/labstack/echo/v4"
)

// Router 路由注册器
type Router struct {
	core  *core.App
	mw    *middleware.Manager
	echo  *echo.Echo
	group *echo.Group
}

// New 创建路由注册器
func New(e *echo.Echo, c *core.App, mw *middleware.Manager) *Router {
	return &Router{
		core:  c,
		mw:    mw,
		echo:  e,
		group: e.Group("/api"),
	}
}

// Register 注册所有路由
func (r *Router) Register() {
	r.echo.GET("/health", func(c echo.Context) error {
		return c.String(200, "OK")
	})

	r.registerPublic()
	r.registerProtected()
}
