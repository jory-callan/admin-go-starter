package router

import (
	"aicode/internal/handler"
)

// registerProtected 注册需要认证的接口
func (r *Router) registerProtected() {
	protected := r.group.Group("")
	protected.Use(r.mw.JWTAuth())

	// 当前用户信息
	protected.GET("/user/info", handler.NewAuthHandler(r.core.DB).GetUserInfo)

	r.registerUserRoutes(protected)
	r.registerRoleRoutes(protected)
	r.registerPermissionRoutes(protected)
}
