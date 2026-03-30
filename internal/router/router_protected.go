package router

import (
	"aicode/internal/handler"
)

// registerProtected 注册需要认证的接口
func (r *Router) registerProtected() {
	// 用户相关接口
	userHandler := handler.NewUserHandler(r.core)
	roleHandler := handler.NewRoleHandler(r.core)

	// 需要 JWT 认证的路由
	auth := r.group.Group("/auth", r.mw.JWTAuth())
	auth.POST("/logout", userHandler.Logout)
	auth.GET("/current_user", userHandler.GetCurrentUser)
	auth.POST("/change_password", userHandler.ChangePassword)

	// 用户管理
	users := r.group.Group("/users", r.mw.JWTAuth())
	users.POST("", userHandler.Create)
	users.GET("", userHandler.List)
	users.GET("/:id", userHandler.GetByID)
	users.PUT("/:id", userHandler.Update)
	users.DELETE("/:id", userHandler.Delete)
	users.POST("/:id/roles", userHandler.AssignRoles)

	// 角色管理
	roles := r.group.Group("/roles", r.mw.JWTAuth())
	roles.POST("", roleHandler.Create)
	roles.GET("", roleHandler.List)
	roles.GET("/all", roleHandler.GetAll)
	roles.GET("/:id", roleHandler.GetByID)
	roles.PUT("/:id", roleHandler.Update)
	roles.DELETE("/:id", roleHandler.Delete)
}
