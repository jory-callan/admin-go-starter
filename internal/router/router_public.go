package router

import (
	"aicode/internal/handler"
)

// registerPublic 注册公开接口（无需认证）
func (r *Router) registerPublic() {
	// 用户公开接口
	userHandler := handler.NewUserHandler(r.core)

	auth := r.group.Group("/auth")
	auth.POST("/login", userHandler.Login)
	auth.POST("/register", userHandler.Register)
}
